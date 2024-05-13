import globals from "globals";
import pluginJs from "@eslint/js";
import pluginVue from "eslint-plugin-vue";
import pluginPrettier from "eslint-plugin-prettier/recommended";

export default [
  { languageOptions: { globals: globals.browser } },
  { ignores: ["dist", "**/wasm_exec.js"] },
  pluginJs.configs.recommended,
  ...pluginVue.configs["flat/recommended"],
  pluginPrettier,
  {
    rules: {
      "no-unused-vars": ["error", { varsIgnorePattern: "^_", argsIgnorePattern: "^_" }],
      "vue/no-template-shadow": "off",
    },
  },
];
