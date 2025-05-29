package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes (server *gin.Engine) {
	// Utility to flush tables without dropping db
	server.GET("/drop_tables", DropTables)

	//		   RECIPES
	// Une route pour déclencher la récupération de recettes via un json contenant un tableau
	server.GET("/recipe", GetRecipe)
	server.POST("/gather", RecipeGathererHandler)

	//         USER
	// Une route pour créer un user + l'enregistrer
	// Une route pour se login
	server.POST("/register", RegisterUser)

	//         PLANNING
	// Une route pour générer un planning + l'enregistrer
	// Une route pour supprimer un planning
	// Une route pour modifier un planning
	// Une route pour changer une recette dans un planning
	server.POST("/planning/generate", GenerateWeeklyPlanning)
	server.POST("/planning/delete", DeletePlanning)
	server.POST("/planning/recipe/change", ChangeRecipeInPlanning)	
	server.GET("/planning", GetPlanning)

	//authenticated := server.Group("/")
	//authenticated.Use(middlewares.Authenticate)
	//authenticated.POST("/events",createEvent)
	//authenticated.PUT("/events/:id", updateEvent)
	//authenticated.DELETE("/events/:id/register", cancelRegistration)
	//server.POST("/login", login)
}