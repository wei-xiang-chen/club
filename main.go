package main

import (
	"club/clients"
	"club/controllers/club"
	"club/controllers/code"
	"club/controllers/login"
	"club/middlewares"
	schedule "club/schedules"
	"club/setting"
	"club/ws/club_ws"
	"club/ws/user_ws"
	"log"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

func init() {
	var err error

	err = setupSetting()
	if err != nil {
		log.Fatalf("inti.setupSetting err : %v", err)
	}

	err = setupDBEngine()
	if err != nil {
		log.Fatalf("inti.setupDBSetting err : %v", err)
	}
}

func main() {
	go club_ws.H.Run()
	go user_ws.H.Run()
	go schedule.DeleteExpiredUser()

	router := initializeRoutes()
	http.ListenAndServe(":8080", router)
}

func initializeRoutes() http.Handler {

	router := gin.Default()
	router.Static("/api-docs", "./swagger/dist")

	v1Router := router.Group("/api/v1/")
	{
		loginRouter := v1Router.Group("/login/").Use(gzip.Gzip(gzip.DefaultCompression))
		{
			loginRouter.POST("/", middlewares.ErrorHandler(login.Login))
		}

		clubRouter := v1Router.Group("/club/").Use(gzip.Gzip(gzip.DefaultCompression)).Use(middlewares.UidAuth())
		{
			clubRouter.GET("/", middlewares.ErrorHandler(club.GetList))
			clubRouter.POST("/", middlewares.ErrorHandler(club.Create))
			clubRouter.POST("/join/:clubId", middlewares.ErrorHandler(club.Join))
			clubRouter.POST("/leave/", middlewares.ErrorHandler(club.Leave))
		}

		codeRouter := v1Router.Group("/code/").Use(gzip.Gzip(gzip.DefaultCompression)).Use(middlewares.UidAuth())
		{
			codeRouter.GET("/", middlewares.ErrorHandler(code.Code))
		}

		wsRouter := v1Router.Group("/ws/")
		{
			wsRouter.GET("/user/:userId", middlewares.ErrorHandler(user_ws.ServeWs))
			wsRouter.GET("/club/:clubId", middlewares.WsErrorHandler(club_ws.ServeWs))
		}
	}

	return cors.AllowAll().Handler(router)
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}

	err = setting.ReadDBSetting()
	if err != nil {
		return err
	}

	return nil
}

func setupDBEngine() error {
	var err error

	clients.DBEngine, err = clients.NewDBEngine(setting.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}
