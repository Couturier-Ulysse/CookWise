import { useMemo } from "react";
import Button from '@mui/material/Button';
import DeleteForeverIcon from '@mui/icons-material/DeleteForever';
import AutoModeIcon from '@mui/icons-material/AutoMode';

const daysOfWeek = [
  "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"
];

export default function UserInputs({ 
  hasRecipes, 
  onGeneratePlanning, 
  onDeletePlanning,
  totalMeals,
  setTotalMeals,
  totalRecipes,
  setTotalRecipes,
  selectedDays = [],
  setSelectedDays
}) {
  // Compute selected days count for totalMeals
  const mealsCount = useMemo(() => selectedDays.length, [selectedDays]);
  const handleDayToggle = (day) => {
    if (hasRecipes) return;
    if (selectedDays.includes(day)) {
      setSelectedDays(selectedDays.filter(d => d !== day));
    } else {
      setSelectedDays([...selectedDays, day]);
    }
  };

  return (
    <div className="grid grid-cols-3 grid-rows-1 gap-3">
      <div className="flex flex-col gap-2">
        {daysOfWeek.map(day => (
          <label key={day} className="flex items-center gap-2">
            <input
              type="checkbox"
              checked={selectedDays.includes(day)}
              onChange={() => handleDayToggle(day)}
              disabled={hasRecipes}
            />
            <span>{day}</span>
          </label>
        ))}
      </div>
      <div>
        <input
          type="number"
          value={totalRecipes}
          min={1}
          max={mealsCount || 7}
          onChange={e => setTotalRecipes(Number(e.target.value))}
          placeholder={hasRecipes ? " " : "Enter number of recipes"}
          className="border rounded-lg p-4 shadow-sm bg-white text-1xl font-semibold text-gray-800"
          disabled={hasRecipes}
        />
      </div>

      {!hasRecipes ? (
        <Button
          variant="contained" 
          endIcon={<AutoModeIcon />}
          onClick={onGeneratePlanning}
          className="self-start h-fit py-2 px-4"
          sx={{
          backgroundColor: '#f59e0B',       // correspond à amber-800
          '&:hover': {
            backgroundColor: '#d97706',     // un peu plus foncé pour le hover (amber-900)
          },
          }}
        >
          Generate Planning
        </Button>
      ) : (
        <Button
          variant="contained"
          endIcon={<DeleteForeverIcon />}
          onClick={onDeletePlanning}
          className="self-start h-fit"
          sx={{
            backgroundColor: '#f59e0B',
            '&:hover': {
              backgroundColor: '#d97706',
            },
          }}
        >
          Delete Planning
        </Button>
      )}
    </div>
  );
}