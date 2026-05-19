package models

import (
	"fmt"
)


type Weekday string
const (
    Monday    Weekday = "monday"
    Tuesday   Weekday = "tuesday"
    Wednesday Weekday = "wednesday"
    Thursday  Weekday = "thursday"
    Friday    Weekday = "friday"
    Saturday  Weekday = "saturday"
    Sunday    Weekday = "sunday"
)

// Struct pour représenter un planning
type Planning struct {
	ID 			 int
    UserID       int
    Year         int
    WeekNumber   int
    RecipesByDay map[Weekday][]Recipe
}



func NewPlanning(user *User, weekNumber int, year int) *Planning {
	return &Planning{
		UserID:			user.ID,
		WeekNumber:		weekNumber,
		Year:			year,
		RecipesByDay: make(map[Weekday][]Recipe),
	}
}

func (p *Planning) AddRecipe(day Weekday, recipe Recipe) {
    p.RecipesByDay[day] = append(p.RecipesByDay[day], recipe)
}

func (p *Planning) SetID(index int) {
	p.ID = index
}

func (p *Planning) PrintInfo() {
	fmt.Printf("🗓️  Planning de la semaine %d (%d)\n", p.WeekNumber, p.Year)
	fmt.Printf("👤 Utilisateur ID : %d\n", p.UserID)
	fmt.Println("📅 Recettes par jour :")

	for _, day := range []Weekday{Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday} {
		fmt.Printf("\n%s:\n", capitalize(string(day)))
		recipes, exists := p.RecipesByDay[day]
		if !exists || len(recipes) == 0 {
			fmt.Println("  (aucune recette)")
			continue
		}
		for i, r := range recipes {
			fmt.Printf("  %d. %s \n", i+1, r.Name)
		}
	}
}
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(s[0]-32) + s[1:]
}
