package service

import (
	"errors"

	logger "technopark-dbms-forum/pkg"

	systemDelivery "technopark-dbms-forum/internal/system/delivery"
	systemRepository "technopark-dbms-forum/internal/system/repository"

	threadDelivery "technopark-dbms-forum/internal/threads/delivery"

	postDelivery "technopark-dbms-forum/internal/posts/delivery"

	forumDelivery "technopark-dbms-forum/internal/forums/delivery"
	userDelivery "technopark-dbms-forum/internal/users/delivery"

	threadRepository "technopark-dbms-forum/internal/threads/repository"

	postRepository "technopark-dbms-forum/internal/posts/repository"
	userRepository "technopark-dbms-forum/internal/users/repository"

	forumRepository "technopark-dbms-forum/internal/forums/repository"

	threadUsecase "technopark-dbms-forum/internal/threads/usecase"

	postUsecase "technopark-dbms-forum/internal/posts/usecase"

	userUsecase "technopark-dbms-forum/internal/users/usecase"

	forumUsecase "technopark-dbms-forum/internal/forums/usecase"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo

	forumUsecase  *forumUsecase.ForumUsecase
	userUsecase   *userUsecase.UserUsecase
	postUsecase   *postUsecase.PostUsecase
	threadUsecase *threadUsecase.ThreadUsecase

	forumRepo  *forumRepository.Postgres
	userRepo   *userRepository.Postgres
	postRepo   *postRepository.Postgres
	threadRepo *threadRepository.Postgres
	systemRepo *systemRepository.Postgres

	forumHandler  *forumDelivery.Handler
	userHanlder   *userDelivery.Handler
	postHandler   *postDelivery.Handler
	threadHandler *threadDelivery.Handler
	systemHandler *systemDelivery.Handler
}

func NewServer(newEcho *echo.Echo) *Server {
	return &Server{
		echo: newEcho,
	}
}

func (s *Server) Start(addr, pgURL string) error {
	if s.echo == nil {
		return errors.New("initialize server first")
	}
	if err := s.init(pgURL); err != nil {
		return errors.New("initialize server error: " + err.Error())
	}

	return s.echo.Start(addr)
}

func (s *Server) init(pgURL string) error {
	if err := s.makeRepositories(pgURL); err != nil {
		return err
	}

	s.makeUseCases()
	s.makeHandlers()
	s.makeRoutes()

	return nil
}

func (s *Server) makeRepositories(url string) (err error) {
	if s.forumRepo, err = forumRepository.NewPostgres(url); err != nil {
		return err
	}
	if s.userRepo, err = userRepository.NewPostgres(url); err != nil {
		return err
	}
	if s.postRepo, err = postRepository.NewPostgres(url); err != nil {
		return err
	}
	if s.threadRepo, err = threadRepository.NewPostgres(url); err != nil {
		return err
	}
	if s.systemRepo, err = systemRepository.NewPostgres(url); err != nil {
		return err
	}

	return nil
}

func (s *Server) makeUseCases() {
	s.forumUsecase = forumUsecase.NewForumUsecase(s.forumRepo)
	s.userUsecase = userUsecase.NewUserUsecase(s.userRepo)
	s.postUsecase = postUsecase.NewPostUsecase(s.postRepo)
	s.threadUsecase = threadUsecase.NewThreadUsecase(s.threadRepo, s.postRepo)
}

func (s *Server) makeHandlers() {
	s.forumHandler = forumDelivery.NewHandler(s.forumUsecase, s.userUsecase)
	s.userHanlder = userDelivery.NewHandler(s.userUsecase)
	s.postHandler = postDelivery.NewHandler(s.postUsecase, s.userUsecase, s.forumUsecase, s.threadUsecase)
	s.threadHandler = threadDelivery.NewHandler(s.threadUsecase, s.forumUsecase)
	s.systemHandler = systemDelivery.NewHandler(s.systemRepo)
}

func (s *Server) makeRoutes() {
	api := s.echo.Group("/api")
	api.Use(logger.Middleware())

	api.POST("/forum/create", s.forumHandler.Create)
	api.GET("/forum/:slug/details", s.forumHandler.GetDetails)
	api.POST("/forum/:slug/create", s.threadHandler.Create)
	api.GET("/forum/:slug/users", s.forumHandler.GetUsers)
	api.GET("/forum/:slug/threads", s.forumHandler.GetThreads)

	api.GET("/post/:id/details", s.postHandler.GetInfo)
	api.POST("/post/:id/details", s.postHandler.Update)

	api.POST("/thread/:slug_or_id/create", s.threadHandler.CreatePosts)
	api.GET("/thread/:slug_or_id/details", s.threadHandler.GetDetails)
	api.POST("/thread/:slug_or_id/details", s.threadHandler.Update)
	api.GET("/thread/:slug_or_id/posts", s.threadHandler.GetPosts)
	api.POST("/thread/:slug_or_id/vote", s.threadHandler.Vote)

	api.POST("/user/:nickname/create", s.userHanlder.Create)
	api.GET("/user/:nickname/profile", s.userHanlder.Get)
	api.POST("/user/:nickname/profile", s.userHanlder.Update)

	api.GET("/service/status", s.systemHandler.GetInfo)
	api.POST("/service/clear", s.systemHandler.Clear)
}
