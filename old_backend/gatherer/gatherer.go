package gatherer

import (
	"fmt"
	"time"
	"ucouturier.com/JOW_scrapper/storage"
)

func RecipeGatherer(ingredientList []string) error {
	fmt.Printf("[GATHER INFO] : Nombre d'ingrédients demandés : %d\n", len(ingredientList))
	for _, ingredient := range ingredientList {
		recipes, err := Search(ingredient, 50)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, recipe := range recipes {
			err = storage.SaveRecipeToDB(&recipe)
			if err != nil {
				return err
			}
		}
		time.Sleep(5 * time.Second)
	}
	return nil
}