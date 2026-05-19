package db

import (
	"database/sql"
	"fmt"
    _ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	fmt.Println("[INFO] Initialisation de la BDD ...")
	var err error
	DB, err = sql.Open("mysql", "backend:backend@tcp(localhost:3306)/cookwise")

	if err != nil {
		fmt.Println(err)
        panic("Could not connect to database.")
    }
 
    DB.SetMaxOpenConns(10)
    DB.SetMaxIdleConns(5)
 
    createTables()
}

func createTables() {

	createRecipeTable :=
	`
	CREATE TABLE IF NOT EXISTS recipes (
		id INT PRIMARY KEY AUTO_INCREMENT,
		query_ID VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		url VARCHAR(255) NOT NULL UNIQUE,
		image_url VARCHAR(255),
		video_url VARCHAR(255),
		description TEXT,
		preparation_time INT NOT NULL,
		cooking_time INT NOT NULL,
		preparation_extra_time_per_cover INT
	)
	`

	_, err := DB.Exec(createRecipeTable)

	if err != nil {
		fmt.Println(err)
        return
    }

    // Table des ingrédients
    createIngredientsTable := `
	CREATE TABLE IF NOT EXISTS ingredients (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL UNIQUE
	)
    `
    _, err = DB.Exec(createIngredientsTable)
    if err != nil {
        fmt.Println("[ERROR] Impossible de créer la table ingredients:", err)
        return
    }

	// Table de liaison recette-ingrédient
	createRecipeIngredientsTable := `
	CREATE TABLE IF NOT EXISTS recipe_ingredients (
		recipe_id INT NOT NULL,
		ingredient_id INT NOT NULL,
		quantity DOUBLE NOT NULL,
		unit VARCHAR(50) NOT NULL,
		FOREIGN KEY (recipe_id) REFERENCES recipes(id),
		FOREIGN KEY (ingredient_id) REFERENCES ingredients(id),
		PRIMARY KEY (recipe_id, ingredient_id)
	)
	`
	_, err = DB.Exec(createRecipeIngredientsTable)
	if err != nil {
		fmt.Println("[ERROR] Impossible de créer la table recipe_ingredients:", err)
		return
	}

	createUsersTable :=
	`
	CREATE TABLE IF NOT EXISTS users (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	)
	`

	_, err = DB.Exec(createUsersTable)

	if err != nil {
		fmt.Println(err)
        return
    }

	// Table de liaison recette-ingrédient
	createPlanningTable := `
	CREATE TABLE IF NOT EXISTS plannings (
		id INT PRIMARY KEY AUTO_INCREMENT,
		user_id INT NOT NULL,
		year INT NOT NULL,
		week_number INT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createPlanningTable)
	if err != nil {
		fmt.Println("[ERROR] Impossible de créer la table plannings:", err)
		return
	}

	createPlanningRecipeTable := `
	CREATE TABLE IF NOT EXISTS planning_recipes (
		id INT PRIMARY KEY AUTO_INCREMENT,
		planning_id INT NOT NULL,
		recipe_id INT NOT NULL,
		weekday ENUM('monday', 'tuesday', 'wednesday', 'thursday', 'friday', 'saturday', 'sunday') NOT NULL,
		FOREIGN KEY (planning_id) REFERENCES plannings(id) ON DELETE CASCADE,
		FOREIGN KEY (recipe_id) REFERENCES recipes(id)
	)
	`
	_, err = DB.Exec(createPlanningRecipeTable)
	if err != nil {
		fmt.Println("[ERROR] Impossible de créer la table planning_recipes:", err)
		return
	}





	fmt.Println("[INFO] Tables créées avec succès")
}