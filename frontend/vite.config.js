import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      "/api": "http://localhost:5555",
      "/thumb": "http://localhost:5555",
      // Only the media-streaming route, not "/content" — that prefix would
      // also swallow the frontend's own "/contents" SPA route.
      "/content/media": "http://localhost:5555",
    },
  },
  build: {
    outDir: "dist",
  },
});
