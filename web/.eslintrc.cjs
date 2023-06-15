/* eslint-env node */
require("@rushstack/eslint-patch/modern-module-resolution");

module.exports = {
  root: true,
  extends: [
    "plugin:vue/vue3-essential",
    "eslint:recommended",
    "@vue/eslint-config-prettier/skip-formatting",
  ],
  rules: {
    "object-curly-spacing": ["error", "always"],
    "require-jsdoc": "off",
    indent: ["error", 2, { SwitchCase: 1 }],
    "no-unused-vars": ["error", { varsIgnorePattern: "^_", argsIgnorePattern: "^_" }],
    "valid-jsdoc": "off",
    "new-cap": "off",
  },
  parserOptions: {
    ecmaVersion: "latest",
  },
  ignorePatterns: ["wasm_exec.js"],
};
