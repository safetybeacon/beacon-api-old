package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/safetybeacon/beacon-api/db"
	_ "github.com/safetybeacon/beacon-api/docs"
	"github.com/safetybeacon/beacon-api/pkg/handlers"
)

var (
	port string = os.Getenv("PORT")
	//
	postgresUser     string = os.Getenv("POSTGRES_USER")
	postgresPassword string = os.Getenv("POSTGRES_PASSWORD")
	postgresDatabase string = os.Getenv("POSTGRES_DATABASE")
	postgresHost     string = os.Getenv("POSTGRES_HOST")
	//
	postgresDSN string = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable",
		postgresUser, postgresPassword,
		postgresHost, postgresDatabase,
	)
)

func init() {
	conn, err := db.NewDB(postgresDSN)
	if err != nil {
		log.Fatalf("failed to connect to the database %v", err)
	}
	defer conn.Close()

	if err := conn.CreateTables(); err != nil {
		log.Fatalf("failed to create database tables %v", err)
	}
}

// @title beacon-api
// @version 1.0
// @description Safety Beacon API

// @host localhost:8080
// @BasePath /v1

func main() {

	router := gin.Default()

	handler := handlers.Handler{
		PostgresLink: postgresDSN,
	}

	v1 := router.Group("/v1")

	swaggerDocs := v1.Group("/swagger")
	{
		swaggerDocs.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	auth := v1.Group("/auth")
	{
		auth.POST("/register", handler.HandleRegister)
		auth.POST("/login", handler.HandleLogin)
		auth.Use(handler.AuthorizeHeader())
		auth.DELETE("/:userid/logout", handler.HandleLogout)
	}

	locations := v1.Group("/locations")
	locations.Use(handler.AuthorizeHeader())
	{
		locations.POST("/:userid", handler.HandleAddLocation)
		locations.GET("/", handler.HandleGetLocations)
	}

	router.Run(fmt.Sprintf(":%s", port))
}
