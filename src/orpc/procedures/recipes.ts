import { os } from "@orpc/server";
import { z } from "zod";

export const recipeRouter = {

  create: os
    .input(
      z.object({
        title: z.string().min(1),
        description: z.string().optional(),
      }),
    )
    .handler(async ({ input }) => {
      return input.title;
    }),

};