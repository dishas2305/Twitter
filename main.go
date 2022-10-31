package main

import (
	"fmt"
	"log"
	"os"

	"twitter/config"
	"twitter/storage"
	route "twitter/v1"

	//_ "twitter/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	//echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		err := godotenv.Load("/var/api/twitter/.env")
		if err != nil {
			log.Fatalf("Error getting env, not comming through %v", err)
		}
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fmt.Println("get env", os.Getenv("ENV"))
	if envName := os.Getenv("ENV"); envName == config.Qa || envName == config.Prod {
		e.Use(middleware.Gzip())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	storage.ConnectLogrus() // log file
	storage.MONGO_DB = storage.ConnectMongoDB()
	fmt.Println("storage.mongo", storage.MONGO_DB)
	//e.GET("/swagger/*", echoSwagger.WrapHandler)
	v1 := e.Group("/api/v1")
	route.InitializeRoutes(v1)
	e.Logger.Fatal(e.Start(":3100"))
}
