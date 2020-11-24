package api

import (
	"app/bot"
	"app/models"
	"app/utils"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	jwtVerify "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type APIServer struct {
	r  *gin.Engine
	db *sqlx.DB
	BotService *bot.Service
	connections map[int]*websocket.Conn
}

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (a *APIServer) GetConnections() map[int]*websocket.Conn {
	return a.connections
}

func (a *APIServer) Start() error {
	a.connections = map[int]*websocket.Conn{}
	var err error
	a.db, err = sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	a.r = gin.Default()
	a.r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Content-Type", "Authorization"},
		AllowMethods:    []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
	}))

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Key:         []byte(os.Getenv("SECRET_KEY")),
		Timeout:     time.Hour,
		IdentityKey: "id",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":         v.Id,
					"country_id": v.CountryId,
					"username":   v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Id:        int(claims["id"].(float64)),
				CountryId: int(claims["country_id"].(float64)),
				Username:  claims["username"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBindJSON(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			username := loginVals.Username
			password := loginVals.Password

			user, err := a.getUser(username)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			if utils.VerifyPassword(user.Password, password) {
				return &models.User{
					Id:        user.Id,
					Username:  user.Username,
					CountryId: user.CountryId,
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":  code,
				"error": message,
			})
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	a.r.POST("/api/user/create", a.handleCreateUser())
	a.r.POST("/api/login", authMiddleware.LoginHandler)
	a.r.GET("/ws", func(c *gin.Context) {
		a.wsHandler(c.Writer, c.Request)
	})

	auth := a.r.Group("/api/auth")
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/questions", a.handleQuestions())
		auth.GET("/questions/:id", a.handleGetQuestionById())
		auth.PATCH("/changeStatus", a.handleChangeStatus())
		auth.GET("/qa", a.handleQAs())
		auth.POST("/qa", a.handleAddQA())
		auth.GET("/qa/:id", a.handleGetQAById())
		auth.POST("/qa/:id", a.handleUpdateQAById())
	}

	return a.r.Run(":8080")
}

func (a *APIServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	wsUpgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer func() {
		for k, _ := range a.connections {
			if a.connections[k] == conn {
				delete(a.connections, k)
			}
		}
		conn.Close()
	}()

	type request struct {
		Type  string      `json:"type" binding:"required"`
		Data  interface{} `json:"data" binding:"required"`
		Token string      `json:"token" binding:"required"`
	}

	var req request

	for {
		err := conn.ReadJSON(&req)
		if err != nil {
			break
		}
		token, err := jwtVerify.Parse(req.Token, func(token *jwtVerify.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwtVerify.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil {
			_ = conn.WriteMessage(1, []byte("0"))
			continue
		}
		if token.Valid {
			if req.Type == "SEND_MESSAGE" {
				mp := req.Data.(map[string]interface{})
				chatId := int64(mp["chatId"].(float64))
				message := mp["message"].(string)
				questionId := int(mp["questionId"].(float64))

				msg := a.sendMessage(chatId, questionId, message)
				a.BotService.SendMessage(chatId, message)

				_ = conn.WriteJSON(map[string]interface{}{
					"type": "message",
					"data": msg,
				})
			} else if req.Type == "GET_MESSAGES" {
				mp := req.Data.(map[string]interface{})
				questionId := int(mp["questionId"].(float64))
				a.connections[questionId] = conn
				messages := a.getMessages(questionId)
				_ = conn.WriteJSON(map[string]interface{}{
					"type": "messages",
					"data": messages,
				})
			}
		} else {
			_ = conn.WriteMessage(1, []byte("0"))
		}
	}
}

func (a *APIServer) handleCreateUser() gin.HandlerFunc {
	type request struct {
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		CountryId int    `json:"country_id" binding:"required"`
	}
	return func(c *gin.Context) {
		var r request
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		id, err := a.createUser(&models.User{Username: r.Username, Password: utils.HashPassword(r.Password), CountryId: r.CountryId})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id": id,
		})
	}
}

func (a *APIServer) handleQuestions() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("id")
		questions, err := a.getQuestionsByCountryId(user.(*models.User).CountryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, questions)
	}
}

func (a *APIServer) handleGetQuestionById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		user, _ := c.Get("id")
		question, err := a.getQuestionById(id, user.(*models.User).CountryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, question)
	}
}

func (a *APIServer) handleChangeStatus() gin.HandlerFunc {
	type request struct {
		Status     int `json:"status"`
		QuestionId int `json:"question_id"`
	}
	return func(c *gin.Context) {
		var r request
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		user, _ := c.Get("id")
		err := a.setStatus(r.Status, r.QuestionId, user.(*models.User).CountryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		if r.Status == 2 {
			chatId := a.getChatIdByQuestionId(r.QuestionId)
			a.BotService.SetUserStatus(chatId, "UNKNOWN")
		}
	}
}

func (a *APIServer) handleAddQA() gin.HandlerFunc {
	type request struct {
		CountryId int    `json:"country_id"`
		Question  string `json:"question"`
		Answer    string `json:"answer"`
	}
	return func(c *gin.Context) {
		var r request
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		err := a.addQA(r.Question, r.Answer, r.CountryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}

func (a *APIServer) handleQAs() gin.HandlerFunc {
	return func(c *gin.Context) {
		countryId, err := strconv.Atoi(c.Query("countryId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		qa, err := a.getQAByCountry(countryId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, qa)
	}
}

func (a *APIServer) handleGetQAById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		qa, err := a.getQAById(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, qa)
	}
}

func (a *APIServer) handleUpdateQAById() gin.HandlerFunc {
	type request struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
	}
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		var r request
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		err = a.updateQA(models.QA{Id: id, Question: r.Question, Answer: r.Answer})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	}
}
