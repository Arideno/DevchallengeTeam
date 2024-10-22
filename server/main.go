package main

import (
	"app/api"
	"app/bot"
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

	botService := &bot.Service{}
	go func() {
		if err := botService.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	apiServer := &api.APIServer{
		BotService: botService,
	}
	botService.ApiServer = apiServer
	if err := apiServer.Start(); err != nil {
		log.Fatal(err)
	}
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