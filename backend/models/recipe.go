package models

import "fmt"

type Recipe struct {
	Query_ID                     string       `json:"id"`
	URL                          string       `json:"url"`
	Name                         string       `json:"name"`
	Ingredients                  []Ingredient `json:"ingredients"`
	ImageURL                     string       `json:"imageUrl"`
	VideoURL                     string       `json:"videoUrl"`
	Description                  string       `json:"description"`
	PreparationTime              int          `json:"preparationTime"`
	CookingTime                  int          `json:"cookingTime"`
	PreparationExtraTimePerCover int          `json:"preparationExtraTimePerCover"`
	CoversCount                  int          `json:"coversCount"`
	// Voir si possible de recup d'autres infos
	//Json                         string       `json:"json"` // JSON brut en string
}

func (r *Recipe) PrintRecipeInfo() {
	fmt.Println(r.Query_ID)
	fmt.Println(r.URL)
	fmt.Println(r.Name)
	fmt.Println(r.Ingredients)
	fmt.Println(r.Description)
	fmt.Println(r.PreparationTime)
	fmt.Println(r.CookingTime)
}