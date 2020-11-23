package api

import (
	"app/models"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"strconv"
)

type APIServer struct {
	r *gin.Engine
	db *sqlx.DB
}

func (a *APIServer) Start() error {
	var err error
	a.db, err = sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	a.r = gin.Default()

	a.r.GET("/questions", a.handleQuestions())
	a.r.GET("/questions/:id", a.handleGetQuestionById())
	a.r.PATCH("/changeStatus", a.handleChangeStatus())
	a.r.GET("/qa", a.handleQAs())
	a.r.POST("/qa", a.handleAddQA())
	a.r.GET("/qa/:id", a.handleGetQAById())
	a.r.POST("/qa/:id", a.handleUpdateQAById())

	return a.r.Run(":8080")
}

func (a *APIServer) handleQuestions() gin.HandlerFunc {
	return func(c *gin.Context) {
		countryId, err := strconv.Atoi(c.Query("countryId"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		questions, err := a.getQuestionsByCountryId(countryId)
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
		question, err := a.getQuestionById(id)
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
		Status int `json:"status"`
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
		err := a.setStatus(r.Status, r.QuestionId)
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

func (a *APIServer) handleAddQA() gin.HandlerFunc {
	type request struct {
		CountryId int `json:"country_id"`
		Question string `json:"question"`
		Answer string `json:"answer"`
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
		Answer string `json:"answer"`
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