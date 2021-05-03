package main

import (
	"club/client"
	"club/controller/club"
	"club/controller/login"
	"club/handler"
	"club/middleware"
	"club/setting"
	"log"

	"github.com/gin-gonic/gin"
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

	router := initializeRoutes()
	router.Run(":8080")
}

func initializeRoutes() *gin.Engine {
	// option := cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{"GET", "POST", "OPTIONS", "DELETE"},
	// 	AllowedHeaders: []string{"X-Request-UID", "Content-Type"},
	// }

	router := gin.Default()
	router.Static("/api-docs", "./swagger/dist")

	v1Router := router.Group("/api/v1/")
	{
		loginRouter := v1Router.Group("/login/").Use(middleware.CORSMiddleware())
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

	return router
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
