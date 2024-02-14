/* eslint-disable */
/* prettier-ignore */
// @ts-nocheck
// noinspection JSUnusedGlobalSymbols
export {}
declare global {
  // Functions exposed by the GoNES WASM app
  const Gones: {
    // Save game and exit
    exit(): void;
    // Save current console state
    saveState(): void;
    // Load console state
    loadState(): void;
  };
}
