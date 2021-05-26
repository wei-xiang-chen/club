package main

import (
	"club/client"
	"club/controller/club"
	"club/controller/code"
	"club/controller/login"
	"club/middleware"
	"club/setting"
	"club/ws"
	"log"
	"net/http"

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
	go ws.H.Run()

	router := initializeRoutes()
	http.ListenAndServe(":8080", router)
}

func initializeRoutes() http.Handler {

	router := gin.Default()
	router.Static("/api-docs", "./swagger/dist")

	v1Router := router.Group("/api/v1/")
	{
		loginRouter := v1Router.Group("/login/")
		{
			loginRouter.POST("/", middleware.ErrorHandler(login.Login))
		}

		clubRouter := v1Router.Group("/club/").Use(middleware.UidAuth())
		{
			clubRouter.GET("/", middleware.ErrorHandler(club.GetList))
			clubRouter.POST("/", middleware.ErrorHandler(club.Create))
			clubRouter.POST("/join/:clubId", middleware.ErrorHandler(club.Join))
			clubRouter.POST("/leave/", middleware.ErrorHandler(club.Leave))
		}

		codeRouter := v1Router.Group("/code/").Use(middleware.UidAuth())
		{
			codeRouter.GET("/", middleware.ErrorHandler(code.Code))
		}

		wsRouter := v1Router.Group("/ws/")
		{
			wsRouter.GET("/:roomId", ws.ServeWs)
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

	client.DBEngine, err = client.NewDBEngine(setting.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}
