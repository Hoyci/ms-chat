import { z } from "zod";

const envSchema = z.object({
  API_URL: z.string().default("http://localhost:80/api/v1"),
});

const parsedEnv = envSchema.safeParse(import.meta.env);

if (!parsedEnv.success) {
  console.error(
    "Invalid configuration in environment variables:",
    parsedEnv.error.format()
  );
  throw new Error("Error in environment configuration");
}

export const config = parsedEnv.data;
