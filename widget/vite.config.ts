import { svelte } from "@sveltejs/vite-plugin-svelte";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [svelte()],
  esbuild: {
    // Strip all legal/license comments (e.g. from lucide-svelte) from JS and CSS
    legalComments: 'none',
  },
  server: {
    port: 5174
  }
});
