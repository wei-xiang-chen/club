package main

import (
	"club/client"
	"club/controller/club"
	"club/controller/login"
	"club/handler"
	"club/middleware"
	"club/setting"
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

	handler := initializeRoutes()
	http.ListenAndServe(":8080", handler)
}

func initializeRoutes() http.Handler {
	router := gin.Default()
	router.Static("/api-docs", "./swagger/dist")

	v1Router := router.Group("/api/v1/")
	{
		loginRouter := v1Router.Group("/login/")
		{
			loginRouter.POST("/", login.Login)
		}

		clubRouter := v1Router.Group("/club/").Use(middleware.CORSMiddleware()).Use(handler.UidAuth())
		{
			clubRouter.GET("/", club.GetList)
			clubRouter.POST("/", club.Create)
			clubRouter.POST("/join/:clubId", club.Join)
			clubRouter.POST("/leave/", club.Leave)
		}
	}

	return cors.Default().Handler(router)
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
