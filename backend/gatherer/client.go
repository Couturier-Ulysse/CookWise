package gatherer

import (
	"encoding/json"
	"fmt"
	"ucouturier.com/JOW_scrapper/utils"
	"ucouturier.com/JOW_scrapper/models"
)

const (
	searchURL = "https://api.jow.fr/public/recipe/quicksearch"
	baseURL   = "https://jow.fr/recipes/"
	staticURL = "https://static.jow.fr/"
)

var (
	optionHeaders = map[string]string{
		"accept":          "*/*",
		"accept-language": "fr,fr-FR;q=0.9,en-US;q=0.8,en;q=0.7",
	}
	postHeaders = map[string]string{
		"accept":          "application/json",
		"accept-language": "fr",
		"content-type":    "application/json",
		"x-jow-withmeta":  "1",
	}
)

// Search calls the Jow API with a query and returns results
func Search(query string, limit int) ([]models.Recipe, error) {
	// OPTION request (required by Jow)
	optionParams := map[string]string{
		"start":              "0",
		"availabilityZoneId": "FR",
		"query":              query,
		"limit":              fmt.Sprint(limit),
	}
	if err := utils.DoOptions(searchURL, optionHeaders, optionParams); err != nil {
		return nil, fmt.Errorf("option request failed: %w", err)
	}

	// POST request to get the recipes
	postParams := optionParams
	body := []byte("{}")
	resp, err := utils.DoPost(searchURL, postHeaders, postParams, body)

	if err != nil {
		return nil, fmt.Errorf("post request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	var parsed struct {
		Data []map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	// Convert to JowResult
	results := make([]models.Recipe, 0, len(parsed.Data))
	for _, r := range parsed.Data {
		res, err := mapToJowResult(r)
		if err != nil {
			continue
		}
		results = append(results, res)
	}

	return results, nil
}

// mapToJowResult converts map data to JowResult
func mapToJowResult(data map[string]interface{}) (models.Recipe, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return models.Recipe{}, err
	}

	var result models.Recipe
	void_string := string(jsonBytes)
	if void_string == ""{
		fmt.Println("[ERROR] With JSON on gathering")
	}

	result.Query_ID = getString(data, "_id")
	result.URL = baseURL + result.Query_ID
	result.Name = getString(data, "title")
	result.ImageURL = staticURL + getString(data, "imageUrl")
	result.VideoURL = staticURL + getString(data, "videoUrl")
	result.Description = getString(data, "description")
	result.PreparationTime = getInt(data, "preparationTime")
	result.CookingTime = getInt(data, "cookingTime")
	result.PreparationExtraTimePerCover = getInt(data, "preparationExtraTimePerCover")
	result.CoversCount = getInt(data, "roundedCoversCount")
	result.Ingredients = parseIngredients(data)

	return result, nil
}

// Helpers to extract typed fields
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	return 0
}

func parseIngredients(data map[string]interface{}) []models.Ingredient {
	constituents, ok := data["constituents"].([]interface{})
	if !ok {
		return nil
	}

	var ingredients []models.Ingredient
	for _, c := range constituents {
		ingMap, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		ing := models.Ingredient{
			Name:       getString(ingMap["ingredient"].(map[string]interface{}), "name"),
			Quantity:   getFloat(ingMap["ingredient"].(map[string]interface{}), "quantityPerCover"),
			IsOptional: ingMap["isOptional"].(bool),
			Unit:       resolveUnit(ingMap),
		}
		ingredients = append(ingredients, ing)
	}
	return ingredients
}

func getFloat(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}

func resolveUnit(ingMap map[string]interface{}) string {
	unit := ingMap["unit"].(map[string]interface{})
	unitID := unit["id"]
	ingredient := ingMap["ingredient"].(map[string]interface{})

	if natural, ok := ingredient["naturalUnit"].(map[string]interface{}); ok {
		if natural["_id"] == unitID {
			return getString(natural, "name")
		}
	}

	if alts, ok := ingredient["alternativeUnits"].([]interface{}); ok {
		for _, a := range alts {
			alt := a.(map[string]interface{})
			if alt["unit"].(map[string]interface{})["_id"] == unitID {
				return getString(alt["unit"].(map[string]interface{}), "name")
			}
		}
	}
	return ""
}