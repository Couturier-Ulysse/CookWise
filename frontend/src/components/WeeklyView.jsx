import { useState } from "react";
import { startOfWeek, addWeeks, format, addDays, parseISO } from "date-fns";
import { enUS } from "date-fns/locale"; // ou fr
import Day from "./Day";

const daysOfWeek = ["Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi", "Dimanche"];
const daysOfWeekLowCase = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"];

export default function WeeklyView({ currentMonday, setCurrentMonday, recipesByDay}) {

  const getWeekFromMonday = (mondayDate) => {
    return Array.from({ length: 7 }, (_, i) =>
      format(addDays(mondayDate, i), "yyyy-MM-dd")
    );
  };

  const week = getWeekFromMonday(currentMonday);
  const startDate = parseISO(week[0]);
  const endDate = parseISO(week[6]);
  const formattedStart = format(startDate, "MMMM dd", { locale: enUS });
  const formattedEnd = format(endDate, "MMMM dd", { locale: enUS });
  const goToPreviousWeek = () => setCurrentMonday((prev) => addWeeks(prev, -1));
  const goToNextWeek = () => setCurrentMonday((prev) => addWeeks(prev, 1));


  return (
    <section className="text-center mt-8">
      <div className="flex justify-center items-center gap-4 mb-2">
        <button onClick={goToPreviousWeek} className="text-2xl px-2">←</button>
        <span className="text-lg font-semibold text-gray-700">
          From {formattedStart} to {formattedEnd}
        </span>
        <button onClick={goToNextWeek} className="text-2xl px-2">→</button>
      </div>
      
      <div className="grid grid-cols-7 gap-4 justify-center">
        {week.map((date, idx) => (
          <Day 
          key={date} 
          dayName={daysOfWeek[idx]} 
          date={date} 
          recipes={recipesByDay[daysOfWeekLowCase[idx]]?.[0]?.name || []} 
          recipesImg={recipesByDay[daysOfWeekLowCase[idx]]?.[0]?.imageUrl || []}
          />

        ))}
      </div>
    </section>
  );
}
