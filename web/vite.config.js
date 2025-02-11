import tailwindcss from "@tailwindcss/vite";
import vue from "@vitejs/plugin-vue";
import { promises as fs } from "node:fs";
import Icons from "unplugin-icons/vite";
import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    rollupOptions: {
      input: {
        main: "index.html",
        game: "game.html",
      },
    },
  },
  plugins: [
    vue(),
    tailwindcss(),
    Icons({
      compiler: "vue3",
      autoInstall: true,
      customCollections: {
        gones: {
          icon: () => fs.readFile("./src/assets/images/icon.svg", "utf-8"),
          heading: () => fs.readFile("./src/assets/images/heading.svg", "utf-8"),
        },
      },
    }),
    VitePWA({
      includeAssets: ["favicon.ico"],
      manifest: {
        name: "GoNES",
        short_name: "GoNES",
        id: "/",
        description: "NES emulator written in Go.",
        theme_color: "#b13939",
        background_color: "#b13939",
        icons: [
          {
            src: "/images/android-chrome-192x192.png",
            sizes: "192x192",
            type: "image/png",
          },
          {
            src: "/images/android-chrome-512x512.png",
            sizes: "512x512",
            type: "image/png",
          },
        ],
      },
      workbox: {
        clientsClaim: true,
        globPatterns: ["**/*{js,css,html,woff2,svg}", "assets/gones*.wasm"],
        maximumFileSizeToCacheInBytes: 15000000,
      },
    }),
  ],
});
