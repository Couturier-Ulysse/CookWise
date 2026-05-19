package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"ucouturier.com/JOW_scrapper/gatherer"
	"ucouturier.com/JOW_scrapper/storage"
)

var statusLock bool

// IngredientRequest représente le format JSON attendu dans le corps de la requête
type IngredientRequest struct {
	Ingredients []string `json:"ingredients"`
}

type RecipeRequest struct {
	Recipe string `json:"recipe"`
}

// RecipeGathererHandler gère la requête HTTP pour collecter des recettes
// Le header doit contenir une liste d'ingrédients renseignés par "ingredients"
func RecipeGathererHandler(context *gin.Context) {
	if statusLock {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": "A Gathering is already processing"})
		return
	}
	statusLock = true

	var req IngredientRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": "Invalid request body"})
		return
	}

	fmt.Println("[INFO] Recettes demandées : ")
	fmt.Println(req)

	if len(req.Ingredients) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": "No recipes given"})
		return
	}

    // Créer un canal pour signaler la fin de la tâche
    done := make(chan string)
	go func() {
		err := gatherer.RecipeGatherer(req.Ingredients)
		if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"[REQUEST ERROR]" : "Can't gather recipes"})
		return
		}
		statusLock = false
		done <- "[INFO] Gathering recipes ended successfully"
	}()

	context.JSON(http.StatusOK, gin.H{"[REQUEST INFO]" : " Gathering launched ..."})

	go func() {
		message := <- done
		// Possible de rajouter à la suite une transmission
		// de l'info ailleurs
		fmt.Println(message)
	}()

}



func GetRecipe(context *gin.Context) {
	id := context.Query("id")
	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": "Missing recipe id"})
		return
	}

	returnedRecipe, err := storage.GetRecipeFromDBByQueryID(id)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"status":   "error",
			"recipe":   nil,
			"message":  "Can't fetch recipe informations",
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"message":  "Recipe fetched successfully",
		"recipe":   returnedRecipe,
	})

}