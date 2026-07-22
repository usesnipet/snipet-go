import { svelte } from "@sveltejs/vite-plugin-svelte";
import { resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vite";

const __dirname = fileURLToPath(new URL('.', import.meta.url));

export default defineConfig({
	plugins: [svelte({
		compilerOptions: {
			css: 'injected',
		},
	})],
	resolve: {
		alias: {
			$lib: resolve(__dirname, 'src/lib'),
		},
	},
	esbuild: {
		legalComments: 'none',
	},
	server: {
		port: 5174,
	},
	build: {
		lib: {
			entry: resolve(__dirname, 'src/embed.ts'),
			name: 'SnipetWidget',
			formats: ['iife'],
			fileName: () => 'snipet-widget.js',
		},
		cssCodeSplit: false,
		rollupOptions: {
			output: {
				inlineDynamicImports: true,
			},
		},
	},
});
