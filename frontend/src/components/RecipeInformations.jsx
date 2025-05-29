import { useEffect, useState } from "react";
import axios from "axios";

export default function RecipeInformations({ recipe, setSelectedRecipe, onInfoClick }) {
  const [recipeDetails, setRecipeDetails] = useState(null);

  useEffect(() => {
    const fetchRecipeDetails = async () => {
      if (!recipe?.id) return;
      try {
        const response = await axios.get('http://localhost:3000/recipe', {
          params: { id: recipe.id }
        });
        setRecipeDetails(response.data?.recipe);
        console.log(response.data?.recipe);
      } catch (error) {
        console.error("Erreur lors de la récupération de la recette :", error);
        setRecipeDetails(null);
      }
    };
    fetchRecipeDetails();
  }, [recipe]);

    
  return (
    <div style={{
      position: "fixed",
      top: 0, left: 0, right: 0, bottom: 0,
      background: "rgba(0,0,0,0.5)",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      zIndex: 1000
    }}>
      <div style={{
        background: "#fff",
        color: "#333",
        padding: "2rem",
        borderRadius: "12px",
        minWidth: "350px",
        maxWidth: "90vw",
        boxShadow: "0 4px 24px rgba(0,0,0,0.2)",
        textAlign: "center",
        position: "relative"
      }}>
        <button
          onClick={() => setSelectedRecipe(null)}
          style={{
            position: "absolute",
            top: "1rem",
            right: "1rem",
            background: "#ff7043",
            color: "#fff",
            border: "none",
            borderRadius: "50%",
            width: "32px",
            height: "32px",
            cursor: "pointer",
            fontWeight: "bold",
            fontSize: "1.2rem"
          }}
          aria-label="Fermer"
        >
          ×
        </button>
        {!recipeDetails ? (
          <h2>Chargement...</h2>
        ) : (
          <>
            <h2 style={{ margin: 0 }}>{recipeDetails.name}</h2>
            {recipeDetails.imageUrl && (
              <img
                src={recipeDetails.imageUrl}
                alt={recipeDetails.name}
                style={{ width: "100%", maxWidth: 300, borderRadius: 8, margin: "1rem 0" }}
              />
            )}
            <p style={{ fontStyle: "italic", color: "#666" }}>{recipeDetails.description}</p>
            <div style={{ margin: "1rem 0" }}>
              <strong>Temps de préparation :</strong> {recipeDetails.preparationTime} min<br />
              <strong>Temps de cuisson :</strong> {recipeDetails.cookingTime} min<br />
              {recipeDetails.preparationExtraTimePerCover > 0 && (
                <span>
                  <strong>Temps supplémentaire par couvert :</strong> {recipeDetails.preparationExtraTimePerCover} min<br />
                </span>
              )}
              <strong>Nombre de couverts :</strong> {recipeDetails.coversCount}
            </div>
            <div style={{ textAlign: "left", margin: "1rem 0" }}>
              <strong>Ingrédients :</strong>
              <ul>
                {recipeDetails.ingredients?.map((ing, idx) => (
                  <li key={idx}>
                    {ing.quantity} {ing.unit} {ing.name}
                    {ing.isOptional && <span style={{ color: "#888" }}> (optionnel)</span>}
                  </li>
                ))}
              </ul>
            </div>
            {recipeDetails.videoUrl && (
              <video
                src={recipeDetails.videoUrl}
                controls
                style={{ width: "100%", maxWidth: 300, borderRadius: 8, margin: "1rem 0" }}
              />
            )}
            <div style={{ marginTop: "1rem" }}>
              <a href={recipeDetails.url} target="_blank" rel="noopener noreferrer" style={{ color: "#ff7043" }}>
                Voir la recette originale sur Jow
              </a>
            </div>
          </>
        )}
      </div>
    </div>
  );
}