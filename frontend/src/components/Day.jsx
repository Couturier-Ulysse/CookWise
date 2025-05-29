export default function Day({ dayName, date, recipes, recipesImg }) {
  const hasRecipe = typeof recipes === "string" && recipes.length > 0;

  return (
    <div className="border rounded-lg p-4 shadow-sm bg-white flex flex-col items-center">
      <h3 className="font-bold text-lg">{dayName}</h3>
      <p className="text-sm text-gray-500">{date}</p>
      {hasRecipe && recipesImg && (
        <img
          src={recipesImg}
          alt={recipes}
          className="w-32 h-32 object-cover rounded-md mt-3 shadow"
        />
      )}
      <div className="mt-4 text-gray-700 text-center">
        {hasRecipe ? (
          <span>{recipes}</span>
        ) : (
          <span className="text-gray-400 italic">Aucune recette</span>
        )}
      </div>
    </div>
  );
}