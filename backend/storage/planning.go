package storage

import (
	"ucouturier.com/JOW_scrapper/models"
	"fmt"
	"ucouturier.com/JOW_scrapper/db"
	"database/sql"
)

func SavePlanningToDB (planning *models.Planning) error {
	planning_id,err := CreatePlanningElement(planning)
	if err != nil {
		fmt.Println("[ERROR] Can't create planning element")
		fmt.Println(err)
		return err
	}

	planning.SetID(planning_id)
	err = MapPlanningToRecipe(planning)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func PlanningVerifyIfExists(planning *models.Planning) (bool,error) {
	query := `
	SELECT id FROM plannings
	WHERE week_number = ?
	AND year = ?
	AND user_id = ?
	`
	rows := db.DB.QueryRow(query,planning.WeekNumber,planning.Year,planning.UserID)
	var element int
	err := rows.Scan(&element)

	if err != nil {
        if err == sql.ErrNoRows {
            // Aucun résultat trouvé, le planning peut être inséré
            return false, nil
        }
        // Une autre erreur s'est produite
        return true, err
	// Vérifier si une ligne a été trouvée
	}
	return true,err
}

func PlanningVerifyIfExistsByAttributes(weekNumber int, year int, user_email string) (bool, error) {
	query := `
	SELECT p.id 
	FROM plannings p
	JOIN users u ON p.user_id = u.id
	WHERE p.week_number = ?
	AND p.year = ?
	AND u.email = ?
	`
	var element int
	err := db.DB.QueryRow(query, weekNumber, year, user_email).Scan(&element)

	if err != nil {
		if err == sql.ErrNoRows {
			// No result found, the planning can be inserted
			return false, nil
		}
		// Another error occurred
		return true, err
	}
	// A row was found
	return true, nil
}


func CreatePlanningElement(planning *models.Planning) (int,error) {
	query := `
	INSERT INTO plannings (
		user_id,
		year,
		week_number
	)
	values (
	?, ?, ?
	)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Erreur avec la requête")
		return 0,err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
						planning.UserID,
						planning.Year,
						planning.WeekNumber)

	if err != nil {
		fmt.Println("[ERROR] Impossible d'insérer la recette")
		fmt.Println(err)
		return 0,err
	}

	index,_ := result.LastInsertId()

	// Si aucun problème, true pour "élément inséré" + nil pas d'erreur
	return int(index),nil
}

func MapPlanningToRecipe(planning *models.Planning) error{

	select_planning_query := `
	SELECT id from plannings where id = ?
	`
	var result_plan int
	err := db.DB.QueryRow(select_planning_query,planning.ID).Scan(&result_plan)

	if err != nil {
		fmt.Println("[ERROR] Can't select planning id from DB")
		return err
	}

	insert_query := `
	INSERT INTO planning_recipes (planning_id,recipe_id,weekday)
	VALUES (?, ?, ?)
	`
	select_recipe_query :=`
	SELECT id from recipes where query_id = ?
	`

	for weekday, recipes := range planning.RecipesByDay {
		fmt.Println("Weekday et recipes :")
		fmt.Println(weekday, recipes)
		for _, recipe := range recipes {
			var result_recipe int
			err := db.DB.QueryRow(select_recipe_query, recipe.Query_ID).Scan(&result_recipe)
			if err != nil {
				fmt.Println("[ERROR] Can't select recipe id from DB")
				return err
			}
			_, err = db.DB.Exec(insert_query, result_plan, result_recipe, weekday)
			if err != nil {
				fmt.Println(err)
				fmt.Println("[ERROR] Can't insert into planning_recipes")
				return err
			}
		}
	}

	return nil
}

func GetPlanningFromDB(email string, year int, weekNumber int) (*models.Planning, error){
	var planning models.Planning
	var user_id int

	select_user_query :=
	`
	SELECT id from users where email = ?
	`
	err := db.DB.QueryRow(select_user_query, email).Scan(&user_id)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Can't select user, does this email exists ?")
		return &planning,err
	}

	select_planning_query := `
	SELECT id, user_id, year, week_number from plannings 
	where user_id = ?
	and year = ?
	and week_number = ?
	`
	err = db.DB.QueryRow(select_planning_query, user_id, year, weekNumber).Scan(&planning.ID, &planning.UserID, &planning.Year, &planning.WeekNumber)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Can't select planning, does it exists ?")
		return &planning,err
	}
	
	select_recipes_query := `
	SELECT pr.weekday, r.query_id, r.name, r.image_url
	FROM planning_recipes pr
	JOIN recipes r ON pr.recipe_id = r.id
	WHERE pr.planning_id = ?
	`

	rows, err := db.DB.Query(select_recipes_query, planning.ID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Can't retrieve recipes for the planning")
		return &planning, err
	}
	defer rows.Close()

	planning.RecipesByDay = make(map[models.Weekday][]models.Recipe)

	for rows.Next() {
		var weekday models.Weekday
		var recipe models.Recipe

		err := rows.Scan(&weekday, &recipe.Query_ID, &recipe.Name, &recipe.ImageURL)
		if err != nil {
			fmt.Println(err)
			fmt.Println("[ERROR] Can't scan recipe data")
			return &planning, err
		}

		planning.RecipesByDay[weekday] = append(planning.RecipesByDay[weekday], recipe)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Error occurred during rows iteration")
		return &planning, err
	}

	return &planning,nil
}

func PlanningDeleteInDB(email string, weekNumber int, year int) (int64, error) {
	// Récupérer l'user_id à partir de l'email
	var userID int
	selectUserQuery := `SELECT id FROM users WHERE email = ?`

	err := db.DB.QueryRow(selectUserQuery, email).Scan(&userID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Can't find user with this email")
		return 0, err
	}

	// Supprimer le planning correspondant
	deleteQuery := `
	DELETE FROM plannings
	WHERE user_id = ?
	AND year = ?
	AND week_number = ?
	`

	stmt, err := db.DB.Prepare(deleteQuery)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Erreur avec la requête")
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(userID, year, weekNumber)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	deletedRows, _ := result.RowsAffected()
	fmt.Println("Lignes supprimées : ", deletedRows)
	return deletedRows, nil
}


func ReplaceRecipeInPlanning(email string, year int, weekNumber int, recipe string) (int64,error) {

	var userID int
	var planningID int

    newRecipeID, err := GetRandomRecipe()

    if err != nil {
        fmt.Println("[ERROR]: can't get a new random recipe")
		fmt.Println(err)
    }

	selectUserQuery := `SELECT id FROM users WHERE email = ?`

	err = db.DB.QueryRow(selectUserQuery, email).Scan(&userID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Can't find user with this email")
		return 0,err
	}

	select_planning_query := `
	SELECT id from plannings
	where user_id = ?
	and year = ?
	and week_number = ?
	`

	err = db.DB.QueryRow(select_planning_query, userID, year, weekNumber).Scan(&planningID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Can't find user with this email")
		return 0,err
	}

	replace_recipes_query := `
	UPDATE planning_recipes
	SET recipe_id = ?
	WHERE planning_id = ?
	  AND recipe_id = (SELECT id FROM recipes WHERE query_id = ?)
	`

	stmt, err := db.DB.Prepare(replace_recipes_query)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Erreur avec la requête")
		return 0,err
	}
	defer stmt.Close()

	result, err := stmt.Exec(newRecipeID,planningID,recipe)

	if err != nil {
		fmt.Println("[ERROR] Impossible d'insérer la recette")
		fmt.Println(err)
		return 0,err
	}

	affectedRows, _ := result.RowsAffected()
	fmt.Println("Lignes modifiées : ", affectedRows)
	return affectedRows, nil
}