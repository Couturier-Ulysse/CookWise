package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"ucouturier.com/JOW_scrapper/models"
	"ucouturier.com/JOW_scrapper/storage"
)

var allWeekdays = []models.Weekday{models.Monday, models.Tuesday, models.Wednesday, models.Thursday, models.Friday, models.Saturday, models.Sunday}

// Struct pour représenter un planning
type GeneratePlanningRequest struct {
    UserEmail     string `json:"user_email"`
    Year          int `json:"year"`
    WeekNumber    int `json:"week_number"`
    TotalMeals    int    `json:"total_meals"`
    TotalRecipes  int    `json:"total_recipes"`
    SelectedDays  []string    `json:"selectedDays"`
}

type DeletePlanningRequest struct {
	UserEmail     string `json:"user_email"`
    Year          int `json:"year"`
    WeekNumber    int `json:"week_number"`
}

type ChangeRecipeInPlanningRequest struct {
   	UserEmail     string `json:"user_email"`
    Year          int `json:"year"`
    WeekNumber    int `json:"week_number"`
    Recipe        string `json:"recipe"`
}


func GenerateWeeklyPlanning (context *gin.Context) {

    var req GeneratePlanningRequest
    if err := context.ShouldBindJSON(&req); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Invalid JSON",
        })
        return
    }

    if req.UserEmail == "" || req.Year == 0 || req.WeekNumber == 0 || req.TotalMeals == 0 || req.TotalRecipes == 0 || len(req.SelectedDays) == 0 {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Missing parameters",
        })
        return
    }

	user,err := storage.GetUserFromDBByEmail(req.UserEmail)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": err})
		return
	}

	planning := models.NewPlanning(user,req.WeekNumber,req.Year)
    recipes, err := storage.GetRandomRecipes(req.TotalMeals,req.TotalRecipes)

	if err != nil {
		context.JSON(http.StatusServiceUnavailable, gin.H{"[REQUEST ERROR]": "Error getting recipes from DB"})
		return
	}

    // Crée une map pour lookup rapide des jours sélectionnés
    selectedDaysMap := make(map[string]bool)
    for _, d := range req.SelectedDays {
        selectedDaysMap[d] = true
    }

    fmt.Println("Recettes piochées : ")
    fmt.Println(recipes)
    fmt.Println("Jours séléctionnés : ")
    fmt.Println(selectedDaysMap)
    fmt.Println("Total meals : ")
    fmt.Println(req.TotalMeals)

    for i,day := range req.SelectedDays {
        // Convert day (string) to models.Weekday
        weekday := models.Weekday(day)
        planning.RecipesByDay[weekday] = append(planning.RecipesByDay[weekday], *recipes[i])
    }

    //// Parcours du lundi au dimanche, mais ne remplit que les jours sélectionnés
    //for _, day := range allWeekdays {
    //    dayStr := string(day)
    //    if !selectedDaysMap[dayStr] {
    //        continue
    //    }
    //    if recipeIndex >= req.TotalMeals {
    //        break
    //    }
    //    planning.RecipesByDay[day] = append(planning.RecipesByDay[day], *recipes[recipeIndex])
    //    recipeIndex++
    //}
    fmt.Println("Recettes par jour :  ")
    fmt.Println(planning.RecipesByDay)

	exists, err := storage.PlanningVerifyIfExists(planning)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"[REQUEST ERROR]": "Can't verify if planning exists"})
		return
	}
	if exists {
		context.JSON(http.StatusBadRequest, gin.H{"[REQUEST ERROR]": "A planning for this user already exists for this week"})
		return
	} else {
		err = storage.SavePlanningToDB(planning)
		if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"[REQUEST ERROR]": "Can't save planning to DB"})
		return
		}
		context.JSON(http.StatusOK, gin.H{"[REQUEST INFO] : Planning saved in DB ": planning})
		return
	}
}

func GetPlanning(context *gin.Context) {
    userEmail := context.Query("user_email")
    yearStr := context.Query("year")
    weekStr := context.Query("week_number")

    if userEmail == "" || yearStr == "" || weekStr == "" {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Missing parameters",
        })
        return
    }
    year, err := strconv.Atoi(yearStr)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Invalid year",
        })
        return
    }
    weekNumber, err := strconv.Atoi(weekStr)
    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Invalid week number",
        })
        return
    }

    exists, _ := storage.UserVerifyIfExistsInDB(userEmail)
    if !exists {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Could not get user, does it exist?",
        })
        return
    }

    exists, _ = storage.PlanningVerifyIfExistsByAttributes(weekNumber, year, userEmail)
    if !exists {
        context.JSON(http.StatusOK, gin.H{
            "status":   "ok",
            "planning": nil,
            "message":  "No planning found for this user and week",
        })
        return
    }

    planning, err := storage.GetPlanningFromDB(userEmail, year, weekNumber)
    if err != nil {
        context.JSON(http.StatusInternalServerError, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  err.Error(),
        })
        return
    }
    context.JSON(http.StatusOK, gin.H{
        "status":   "ok",
        "planning": planning,
    })
}

func DeletePlanning(context *gin.Context) {

    var req DeletePlanningRequest
    if err := context.ShouldBindJSON(&req); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Invalid JSON",
        })
        return
    }

	if req.UserEmail == "" || req.Year == 0 || req.WeekNumber == 0 {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Missing parameters",
        })
        return
    }

    exists, _ := storage.UserVerifyIfExistsInDB(req.UserEmail)
    if !exists {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
			"deleted": 0,
            "message":  "Could not get user, does it exist?",
        })
        return
    }

    exists, _ = storage.PlanningVerifyIfExistsByAttributes(req.WeekNumber, req.Year, req.UserEmail)
    if !exists {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
			"deleted": 0,
            "message":  "No planning to delete for this week, verify your delete command",
        })
        return
    } else {

		deletedRows,err := storage.PlanningDeleteInDB(req.UserEmail, req.WeekNumber, req.Year)
		if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
			"deleted": 0,
            "message":  "Can't delete planning",
        }) 
		return
		}
		context.JSON(http.StatusOK, gin.H{
            "status":   "error",
            "deleted": deletedRows,
			"message":  "PlanningDeleted",
        })
        return
	}


}

func ChangeRecipeInPlanning(context *gin.Context) {

    var req ChangeRecipeInPlanningRequest
    if err := context.ShouldBindJSON(&req); err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Invalid JSON",
        })
        return
    }

	if req.UserEmail == "" || req.Year == 0 || req.WeekNumber == 0 || req.Recipe == "" {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Missing parameters",
        })
        return
    }

    exists, _ := storage.UserVerifyIfExistsInDB(req.UserEmail)
    if !exists {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
			"deleted": 0,
            "message":  "Could not get user, does it exist?",
        })
        return
    }

    modifiedRecipes,err := storage.ReplaceRecipeInPlanning(req.UserEmail, req.Year, req.WeekNumber, req.Recipe)

    if err != nil {
        context.JSON(http.StatusBadRequest, gin.H{
            "status":   "error",
            "planning": nil,
            "message":  "Can't replace recipe in this planning",
        })
        return
    }

    context.JSON(http.StatusOK, gin.H{
        "status":   "error",
        "modified": modifiedRecipes,
        "message":  "Recipes modified successfully",
    })
    return

}
