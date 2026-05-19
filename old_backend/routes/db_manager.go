package routes

import (
	"net/http"
	"github.com/gin-gonic/gin"
    "ucouturier.com/JOW_scrapper/db"
	"fmt"
)

func DropTables(context *gin.Context) {

	dropPlanningsTable :=
	`
	DELETE FROM plannings
	`
	_, err := db.DB.Exec(dropPlanningsTable)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"[REQUEST ERROR]": "Can't drop planning table"})
		return
	}
	dropUsersTable :=
	`
	DELETE FROM users
	`
	_, err = db.DB.Exec(dropUsersTable)
	if err != nil {
		fmt.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"[REQUEST ERROR]": "Can't drop users table"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"[REQUEST INFO]" : "Tables reset"})
	return
}