package storage

import (
	"ucouturier.com/JOW_scrapper/models"
	"ucouturier.com/JOW_scrapper/db"
	"database/sql"
	"fmt"
	"errors"
)


func CreateIngredientElement(ingredient *models.Ingredient) error {
	query := `
	INSERT INTO ingredients (name)
	values (?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		fmt.Println("[ERROR] Erreur avec la requête")
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(ingredient.Name)

	if err != nil {
		fmt.Println("[ERROR] Impossible d'insérer l'ingredient")
		return err
	}
	_, _ = result.RowsAffected()

	// Si aucun problème, true pour "élément inséré" + nil pas d'erreur
	return nil
}

func IngredientVerifyIfExists(ingredient *models.Ingredient) (bool,error) {
	query := `
	SELECT name FROM ingredients
	WHERE name = ?
	`
	rows := db.DB.QueryRow(query,ingredient.Name)
	var element string
	err := rows.Scan(&element)

	if err != nil {
        if err == sql.ErrNoRows {
            // Aucun résultat trouvé, l'ingrédient peut être inséré
            return false, nil
        }
        // Une autre erreur s'est produite
		fmt.Println(err)
        return true, errors.New("[ERROR] Impossible de lire l'ingrédient")
    }

    // Si l'élément correspond, il existe déjà
    if element == ingredient.Name {
        return true, nil
    }
    // Par défaut, l'élément n'existe pas
    return false, nil
}

func GetIngredientFromDB(name string) (*models.Ingredient, error) {
	query := `SELECT name FROM ingredients WHERE name = ?`

	row := db.DB.QueryRow(query, name)

	var ing models.Ingredient
	err := row.Scan(&ing.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ingredient '%s' not found", name)
		}
		return nil, err
	}

	return &ing, nil
}