import { useState, useEffect } from "react";
import axios from "axios";
import Header from './components/Header.jsx';
import WeeklyView from "./components/WeeklyView";
import { startOfWeek, addWeeks, format, addDays, parseISO } from "date-fns";
import RecipeList from "./components/RecipeList";
import UserInputs from "./components/UserInputs.jsx";
import RecipeInformations from "./components/RecipeInformations.jsx";

// Fonction pour obtenir le numéro de la semaine courante
function getWeekNumber(date = new Date()) {
  const firstDayOfYear = new Date(date.getFullYear(), 0, 1);
  const pastDaysOfYear = (date - firstDayOfYear) / 86400000;
  // 0 = Sunday, 1 = Monday, etc.
  // ISO week: weeks start on Monday, week 1 is the week with the first Thursday of the year
  const dayOfWeek = firstDayOfYear.getDay() || 7;
  const weekNumber = Math.ceil((pastDaysOfYear + dayOfWeek - 1) / 7);
  return weekNumber;
}

export default function App() {
  const [planningData, setPlanningData] = useState(null);
  const User = 'Ulysse';
  const [currentMonday, setCurrentMonday] = useState(
    startOfWeek(new Date(), { weekStartsOn: 1 })
  );
  const [weekNumber, setWeekNumber] = useState(getWeekNumber(currentMonday));
  const [totalMeals, setTotalMeals] = useState(5);
  const [totalRecipes, setTotalRecipes] = useState(3);
  const [selectedDays, setSelectedDays] = useState([]);
  const [selectedRecipe, setSelectedRecipe] = useState(null);

  // Requete pour fetch un planning existant
  const fetchPlanning = async () => {
    try {
      const response = await axios.get('http://localhost:3000/planning', {
        params: {
          user_email: "ulysse@gmail.com",
          year: currentMonday.getFullYear(),
          week_number: getWeekNumber(currentMonday)
        }
      });
      setPlanningData(response.data);
    } catch (error) {
      console.error("Erreur lors de la récupération du planning :", error);
    }
  };

  // Requete pour generer un planning
  const generatePlanning = async () => {
    try {
      await axios.post('http://localhost:3000/planning/generate', {
        user_email: "ulysse@gmail.com",
        year: currentMonday.getFullYear(),
        week_number: getWeekNumber(currentMonday),
        total_meals: selectedDays.length,
        total_recipes: totalRecipes,
        selectedDays: selectedDays
      });
      // Recharge le planning après génération
      fetchPlanning();
    } catch (error) {
      console.error("Erreur lors de la génération du planning :", error);
    }
  };
  const deletePlanning = async () => {
    try {
      await axios.post('http://localhost:3000/planning/delete', {
        user_email: "ulysse@gmail.com",
        year: currentMonday.getFullYear(),
        week_number: getWeekNumber(currentMonday)
      });
      // Recharge le planning après génération
      fetchPlanning();
    } catch (error) {
      console.error("Erreur lors de la suppression du planning :", error);
    }
  }
  const changeRecipeInPlanning = async (recipe) => {
    try {
      await axios.post('http://localhost:3000/planning/recipe/change', {
        user_email: "ulysse@gmail.com",
        year: currentMonday.getFullYear(),
        week_number: getWeekNumber(currentMonday),
        recipe: recipe.id
      });
      // Recharge le planning après génération
      fetchPlanning();
    } catch (error) {
      console.error("Erreur lors de la suppression du planning :", error);
    }
  }

  useEffect(() => {
    setWeekNumber(getWeekNumber(currentMonday));
  }, [currentMonday]);
  useEffect(() => {
    fetchPlanning();
  }, [currentMonday]);


  // Récupère les recettes de la semaine actuelle et vérifie s'il y en a
  const recipesByDay = planningData?.planning?.RecipesByDay || {};
  const hasRecipes = Object.values(recipesByDay).flat().length > 0;
  const recipesByDayObj = planningData?.["planning"]?.["RecipesByDay"] || {};
  const allRecipes = Object.values(recipesByDayObj)
    .flat()
    .filter((r, i, arr) => r && arr.findIndex(x => x.id === r.id) === i);

  return (
    <>
<div
  className={
    "flex flex-col items-center justify-center gap-15 pt-0 " +
    (selectedRecipe ? " pointer-events-none opacity-50" : "")
  }
>
      <Header />
      <UserInputs 
        hasRecipes={hasRecipes} 
        onGeneratePlanning={generatePlanning} 
        onDeletePlanning={deletePlanning}
        totalMeals={totalMeals}
        setTotalMeals={setTotalMeals}
        totalRecipes={totalRecipes}
        setTotalRecipes={setTotalRecipes} 
        selectedDays={selectedDays}
        setSelectedDays={setSelectedDays}
      />
      <WeeklyView
        currentMonday={currentMonday}
        setCurrentMonday={setCurrentMonday}
        recipesByDay={recipesByDayObj}
      />
    </div>

    <div className="w-full flex justify-start px-8">
      <RecipeList
        recipes={allRecipes}
        changeRecipeInPlanning={changeRecipeInPlanning}
        setSelectedRecipe={setSelectedRecipe}
      />
    </div>

  {selectedRecipe !== null && (
    <RecipeInformations recipe={selectedRecipe} 
    setSelectedRecipe={setSelectedRecipe} />
  )}

  </>
  );
}