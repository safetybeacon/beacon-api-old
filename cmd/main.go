package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/t4ke0/locations_api/db"
	"github.com/t4ke0/locations_api/pkg/handlers"
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

func main() {

	router := gin.Default()

	handler := handlers.Handler{
		PostgresLink: postgresDSN,
	}

	v1 := router.Group("/v1")

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
