export default function RecipeList({ recipes, changeRecipeInPlanning, setSelectedRecipe}) {
  return (
    <div className="w-full max-w-2xl mt-8">
      <h2 className="text-xl font-bold mb-2">Recettes de la semaine</h2>
      <ul className="list-none pl-0 flex gap-4">
        {recipes.length === 0 ? (
          <li className="text-gray-400 italic">Aucune recette</li>
        ) : (
          recipes.map((recipe) => (
            <li key={recipe.id} className="flex items-center  mb-2">
              <button
                className="mr-3 text-1xl font-semibold text-amber-800 border rounded-lg p-2 shadow-sm bg-white"
                onClick={() => changeRecipeInPlanning(recipe)}
                aria-label={`Supprimer ${recipe.name}`}
              >
                X
              </button>
              <span className="max-w-xs truncate block">{recipe.name}</span>

              <button
                className="mr-3 text-1xl font-semibold text-blue-400 border rounded-lg p-2 shadow-sm bg-white"
                onClick={() => setSelectedRecipe(recipe)}
              >
                I
              </button>
            </li>
          ))
        )}
      </ul>
    </div>
  );
}