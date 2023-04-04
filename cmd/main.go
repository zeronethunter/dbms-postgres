package main

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"os"
	service "technopark-dbms-forum/internal/init"
	logger "technopark-dbms-forum/pkg"
)

func main() {
	l := logger.GetInstance()

	e := echo.New()

	prometheusEcho := echo.New()
	p := prometheus.NewPrometheus("echo", nil)
	e.Use(p.HandlerFunc)
	p.SetMetricsPath(prometheusEcho)

	e.Logger = l
	prometheusEcho.Logger = l

	s := service.NewServer(e)

	pgURL := "host=" + os.Getenv("DB_HOST") + " port=" + os.Getenv("DB_PORT") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=forum-task sslmode=disable"

	go func() { prometheusEcho.Logger.Fatal(prometheusEcho.Start(":" + os.Getenv("METRICS_PORT"))) }()

	if err := s.Start(":"+os.Getenv("PORT"), pgURL); err != nil {
		log.Error(err)
	}
}
