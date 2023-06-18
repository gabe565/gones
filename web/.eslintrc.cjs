/* eslint-env node */
require("@rushstack/eslint-patch/modern-module-resolution");

module.exports = {
  root: true,
  env: {
    browser: true,
  },
  extends: [
    "plugin:vue/vue3-recommended",
    "eslint:recommended",
    "@vue/eslint-config-prettier/skip-formatting",
    "prettier",
    "plugin:prettier/recommended",
  ],
  rules: {
    "no-unused-vars": ["error", { varsIgnorePattern: "^_", argsIgnorePattern: "^_" }],
    "vue/no-template-shadow": "off",
  },
  parserOptions: {
    sourceType: "module",
  },
  ignorePatterns: ["wasm_exec.js"],
};
