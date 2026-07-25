import { toast } from "@/hooks/use-toast";
import { create } from "zustand";

import { apiKeyKeySchema } from "./schemas";

const KEY = "snipet@api-key";

type ApiKeyStore = {
  key: string | null;
  set: (key: string | null) => void;
};

export const useApiKeyStore = create<ApiKeyStore>((set) => ({
  key: sessionStorage.getItem(KEY),
  set: (key: string | null) => {
    if (key) {
      const parsed = apiKeyKeySchema.safeParse(key);
      if (!parsed.success) {
        toast({
          title: "Invalid API key",
          description: "The API key is invalid",
          variant: "destructive",
        });
        return;
      }
      sessionStorage.setItem(KEY, key);
    } else {
      sessionStorage.removeItem(KEY);
    }

    set({ key });
  },
}));
