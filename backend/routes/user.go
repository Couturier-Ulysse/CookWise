package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"ucouturier.com/JOW_scrapper/storage"
	"ucouturier.com/JOW_scrapper/models"
)

// RecipeGathererHandler gère la requête HTTP pour collecter des recettes
// Le header doit contenir une liste d'ingrédients renseignés par "ingredients"
func RegisterUser(context *gin.Context) {
	var userRequest models.User
	err := context.ShouldBindBodyWithJSON(&userRequest)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": "Could not parse data"})
		return
	}
	
	userExists,err := storage.UserVerifyIfExistsInDB(userRequest.Email)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"[REQUEST ERROR]": "Can't verify if user exists"})
		return
	}

	if userExists {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": "Can't create user : email already exists"})
		return
	} else {
		err := storage.SaveUserToDB(&userRequest)
		if err != nil {
			context.JSON(http.StatusServiceUnavailable, gin.H{"[REQUEST ERROR]": "Error saving User in DB"})
			return
		}
		context.JSON(http.StatusOK, gin.H{"[REQUEST INFO]": "User saved successfully"})
		return
	}

}
