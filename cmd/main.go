package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	service "technopark-dbms-forum/internal/init"
	logger "technopark-dbms-forum/pkg"
)

func main() {
	l := logger.GetInstance()

	e := echo.New()
	e.Logger = l

	s := service.NewServer(e)

	pgURL := "host=localhost port=5432 user=zenehu password=zenehu dbname=forum-task sslmode=disable"

	if err := s.Start(":5000", pgURL); err != nil {
		log.Error(err)
	}
}
