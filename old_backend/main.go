package main

import (
    "ucouturier.com/JOW_scrapper/db"
    "github.com/gin-gonic/gin"
    "ucouturier.com/JOW_scrapper/routes"
    "github.com/gin-contrib/cors"  // ← ajouter cet import
    "time"
)

func main() {
    db.InitDB()
    server := gin.Default()

    // Middleware CORS
    server.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},  // Ton frontend React
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))

    routes.RegisterRoutes(server)
    server.Run("0.0.0.0:3000")
}