import "./wasm_exec";
import wasmUrl from "../assets/gones.wasm?url";
import { waitForElement } from "../util/element";
import { dbGet, dbPut } from "../plugins/db";

// Polyfill
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

// eslint-disable-next-line no-undef
const go = new Go();
let inst;
WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject).then((result) => {
  inst = result.instance;
  // Notify parent that iframe is ready
  window.parent.postMessage({ type: "ready" });
});

// Begin a game
window.addEventListener("message", async ({ data }) => {
  if (data.type === "play") {
    window.cartridge = data.cartridge;
    go.run(inst);
    (await waitForElement("canvas")).focus();
  }
});

// Focus the canvas when iframe is focused
window.addEventListener("focus", () => document.querySelector("canvas")?.focus());

window.GonesClient = {
  SetRomName(value) {
    window.parent.postMessage({ type: "name", value });
  },
  DbPut: dbPut,
  DbGet: dbGet,
};
