import { toast } from "@/hooks/use-toast";
import { create } from "zustand";

import { jwtSchema } from "./schemas";

type JwtStore = {
  token: string | null;
  setToken(key: string | null): void;
};

const KEY = "snipet@access-token";
export const jwtStore = create<JwtStore>((set) => ({
  token: localStorage.getItem(KEY) as string | null,

  setToken: (key) => {
    if (key) {
      const parsed = jwtSchema.safeParse(key);
      if (!parsed.success) {
        toast({
          title: "Invalid API key",
          description: "The API key is invalid",
          variant: "destructive",
        });
        return;
      }
      localStorage.setItem(KEY, key);
    } else {
      localStorage.removeItem(KEY);
    }

    set({ token: key });
  },
}));