package storage

import (
	"ucouturier.com/JOW_scrapper/models"
	"ucouturier.com/JOW_scrapper/db"
	"database/sql"
	"fmt"
	"errors"
)


func SaveRecipeToDB (recipe *models.Recipe) error {

	exists, err := RecipeVerifyIfExists(recipe)

	if err != nil {
		fmt.Println(err)
		return err
	}
	if exists {
		fmt.Println("[STORAGE WARN] " + recipe.Name +  " déjà téléchargée ")
		return nil
	// if no :
	} else {
		// La fonction crée la recette
		err = CreateRecipeElement(recipe)
		if err != nil {
			fmt.Println(err)
			return err
		}

		for i := 0; i < len(recipe.Ingredients); i++ {
			act_ingredient := recipe.Ingredients[i]
			exists, err = IngredientVerifyIfExists(&act_ingredient)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if !exists {
				err = CreateIngredientElement(&act_ingredient)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
			err = MapIngredientToRecipe(recipe,&act_ingredient)
			if err != nil {
				return err
			}
		}
		fmt.Println("[STORAGE INFO] " + recipe.Name + " ajoutée avec succès")
	}
	return nil

}

func RecipeVerifyIfExists(recipe *models.Recipe) (bool,error) {
	query := `
	SELECT query_ID FROM recipes
	WHERE query_ID = ?
	`
	rows := db.DB.QueryRow(query,recipe.Query_ID)
	var element string
	err := rows.Scan(&element)

	if err != nil {
        if err == sql.ErrNoRows {
            // Aucun résultat trouvé, l'ingrédient peut être inséré
            return false, nil
        }
        // Une autre erreur s'est produite
        return true, errors.New("[ERROR] Impossible de lire l'ingrédient")
    }

    // Si l'élément correspond, il existe déjà
    if element == recipe.Query_ID {
        return true, nil
    }
    // Par défaut, l'élément n'existe pas
    return false, nil
}

func CreateRecipeElement(recipe *models.Recipe) error {
	query := `
	INSERT INTO recipes (
		query_ID ,
		name,
		url,
		image_url,
		video_url,
		description,
		preparation_time,
		cooking_time,
		preparation_extra_time_per_cover
	)
	values (
	?, ?, ?, ?, ?, ?, ?, ?, ?
	)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Erreur avec la requête")
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
						recipe.Query_ID,
						recipe.Name,
						recipe.URL,
						recipe.ImageURL,
						recipe.VideoURL,
						recipe.Description,
						recipe.PreparationTime,
						recipe.CookingTime,
						recipe.PreparationExtraTimePerCover)

	if err != nil {
		fmt.Println("[ERROR] Impossible d'insérer la recette")
		return err
	}
	_, _ = result.RowsAffected()

	// Si aucun problème, true pour "élément inséré" + nil pas d'erreur
	return nil
}

func MapIngredientToRecipe(recipe *models.Recipe, ingredient *models.Ingredient) error{
	var result_ing, result_recipe string
	// Recuperer les id de recette et d'ingredient
	ingredient_query := `SELECT id from ingredients where name = ?`
	err := db.DB.QueryRow(ingredient_query,ingredient.Name).Scan(&result_ing)
	if err != nil {
		fmt.Println("[ERROR] Erreur avec la requête")
		return err
	}
	recipe_query := `SELECT id FROM recipes where query_id = ?`
	err = db.DB.QueryRow(recipe_query,recipe.Query_ID).Scan(&result_recipe)
	if err != nil {
		fmt.Println("[ERROR] Erreur avec la requête")
		return err
	}

	end_query := `
	INSERT INTO recipe_ingredients (
		recipe_id,
		ingredient_id,
		quantity,
		unit
	)
	values (?, ?, ?, ?)
	`
	stmt, err := db.DB.Prepare(end_query)
	if err != nil {
		fmt.Println("[ERROR] Erreur avec la requête")
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
						result_recipe,
						result_ing,
						ingredient.Quantity,
						ingredient.Unit)
	if err != nil {
		fmt.Println("[ERROR] Impossible d'insérer la recette")
		return err
	}
	_, _ = result.RowsAffected()

	return nil
}

func GetRecipeFromDB(id int) (*models.Recipe, error) {
    query := `SELECT query_ID, name, url, image_url, video_url, description, preparation_time, cooking_time, preparation_extra_time_per_cover FROM recipes WHERE id = ?`

    row := db.DB.QueryRow(query, id)

    var recipe models.Recipe
    err := row.Scan(
        &recipe.Query_ID,
        &recipe.Name,
        &recipe.URL,
        &recipe.ImageURL,
        &recipe.VideoURL,
        &recipe.Description,
        &recipe.PreparationTime,
        &recipe.CookingTime,
        &recipe.PreparationExtraTimePerCover,
    )

    if err != nil {
        return nil, err
    }

    return &recipe, nil
}

func GetRecipeFromDBByQueryID(queryID string) (*models.Recipe, error) {

	recipeQuery := `
	SELECT 
	    Query_ID,
        name,
        url,
        image_url,
        video_url,
        description,
        preparation_time,
        cooking_time,
        preparation_extra_time_per_cover
	from recipes
	where query_ID = ?
	`

    row := db.DB.QueryRow(recipeQuery, queryID)
    var recipe models.Recipe
	err := row.Scan(
        &recipe.Query_ID,
        &recipe.Name,
        &recipe.URL,
        &recipe.ImageURL,
        &recipe.VideoURL,
        &recipe.Description,
        &recipe.PreparationTime,
        &recipe.CookingTime,
        &recipe.PreparationExtraTimePerCover,
    )

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Récupérer l'ID de la recette à partir de query_ID
	var recipeID int
	idQuery := `SELECT id FROM recipes WHERE query_ID = ?`
	err = db.DB.QueryRow(idQuery, queryID).Scan(&recipeID)
	if err != nil {
		return nil, err
	}

	// Récupérer tous les ingrédients associés à cette recette
	ingredientQuery := `
		SELECT i.name, ri.quantity, ri.unit
		FROM recipe_ingredients ri
		JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = ?
	`
	rows, err := db.DB.Query(ingredientQuery, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ingredients []models.Ingredient
	for rows.Next() {
		var ing models.Ingredient
		err := rows.Scan(&ing.Name, &ing.Quantity, &ing.Unit)
		if err != nil {
			return nil, err
		}
		ingredients = append(ingredients, ing)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	recipe.Ingredients = ingredients

	fmt.Println(recipe)

	return &recipe, nil
}

func GetRandomRecipes(totalMeals int, totalRecipes int) ([]*models.Recipe, error) {
	query := `
	SELECT query_ID, name, url, image_url, video_url, description, preparation_time, cooking_time, preparation_extra_time_per_cover
	FROM recipes
	ORDER BY RAND()
	LIMIT ?
	`

	rows, err := db.DB.Query(query, totalRecipes)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Erreur avec la requête")
		return nil, err
	}
	defer rows.Close()

	var recipes []*models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.Query_ID,
			&recipe.Name,
			&recipe.URL,
			&recipe.ImageURL,
			&recipe.VideoURL,
			&recipe.Description,
			&recipe.PreparationTime,
			&recipe.CookingTime,
			&recipe.PreparationExtraTimePerCover,
		)
		if err != nil {
			fmt.Println("[ERROR] Erreur lors de la lecture des résultats")
			return nil, err
		}
		recipes = append(recipes, &recipe)
	}

	if err = rows.Err(); err != nil {
		fmt.Println("[ERROR] Erreur lors de l'itération des résultats")
		return nil, err
	}

	var distributedRecipes []*models.Recipe
	mealsPerRecipe := totalMeals / len(recipes)
	remaining := totalMeals % len(recipes)

	for i, recipe := range recipes {
		repeat := mealsPerRecipe
		if i < remaining {
			// Répartit équitablement les repas restants
			repeat += 1
		}
		for j := 0; j < repeat; j++ {
			distributedRecipes = append(distributedRecipes, recipe)
		}
	}


	return distributedRecipes, nil
}

func GetRandomRecipe() (int, error) {
	query := `
	SELECT id
	FROM recipes
	ORDER BY RAND()
	LIMIT 1
	`

	var recipeID int

	err := db.DB.QueryRow(query).Scan(&recipeID)
	if err != nil {
		fmt.Println(err)
		fmt.Println("[ERROR] Erreur avec la requête")
		return 0, err
	}

	if err != nil {
		fmt.Println("[ERROR] Erreur lors de la lecture des résultats")
		return 0, err
	}

	fmt.Println(recipeID)

	return recipeID, nil
}