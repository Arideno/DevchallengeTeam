package main

import (
	"app/bot"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"os"
)

func main() {
	if err := migration(); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := bot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "test",
		})
	})
	log.Fatal(r.Run())
}

func migration() error {
	m, err := migrate.New("file://./migrations", os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}

	return nil
}