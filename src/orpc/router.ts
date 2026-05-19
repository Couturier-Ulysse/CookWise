import { recipeRouter } from "@/orpc/procedures";

export const appRouter = {
  recipes: recipeRouter,
};

export type AppRouter = typeof appRouter;