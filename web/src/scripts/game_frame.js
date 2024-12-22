import wasmUrl from "../assets/gones.wasm?url";
import { dbGet, dbPut } from "../plugins/db";
import { waitForElement } from "../util/element";
import {
  exitEvent,
  loadStateEvent,
  newExitEvent,
  newNameEvent,
  newReadyEvent,
  playEvent,
  saveStateEvent,
} from "../util/events";
import "./wasm_exec";

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
  window.parent.dispatchEvent(newReadyEvent());
});

const handleFocus = async () => {
  const el = await waitForElement("canvas");
  el.focus();
};
window.addEventListener("focus", handleFocus);

window.addEventListener(playEvent, async ({ detail: { name, data } }) => {
  window.GonesCartridge = { name, data: new Uint8Array(await data) };
  await Promise.all([go.run(inst), handleFocus()]);
  window.parent.dispatchEvent(newExitEvent());
});

window.addEventListener(exitEvent, () => {
  if (window.Gones) window.Gones.exit();
});

window.addEventListener(saveStateEvent, () => {
  if (window.Gones) window.Gones.saveState();
});

window.addEventListener(loadStateEvent, () => {
  if (window.Gones) window.Gones.loadState();
});

window.GonesClient = Object.freeze({
  setRomName(value) {
    window.parent.dispatchEvent(newNameEvent(value));
  },
  dbPut,
  dbGet,
});
