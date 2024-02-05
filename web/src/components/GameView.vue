<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import MenuButton from "./MenuButton.vue";
import SettingsMenu from "./SettingsMenu.vue";
import { wait } from "../util/wait";

const showSettings = ref(true);

// Promise that will resolve when the iframe is done reloading
let resolve;
let promise = new Promise((r) => {
  resolve = r;
});

// Handle messages from the iframe
const iframeMessage = ({ data }) => {
  if (data.type === "ready") {
    resolve();
  }
};
onMounted(() => {
  window.addEventListener("message", iframeMessage);
});
onBeforeUnmount(() => {
  window.removeEventListener("message", iframeMessage);
});

const iframe = ref();
let running = false;

const cartridgeInserted = async (val) => {
  showSettings.value = false;
  if (running) {
    running = true;
    promise = new Promise((r) => {
      resolve = r;
    });
    iframe.value.contentWindow.Gones.exit();
    await wait(100);
    await iframe.value.contentWindow.location.reload();
    await promise;
  }
  iframe.value.contentWindow.postMessage({
    type: "play",
    cartridge: new Uint8Array(await val.arrayBuffer()),
  });
  iframe.value.contentWindow.focus();
  running = true;
};

watch(showSettings, (val) => {
  if (!val && running) {
    iframe.value.contentWindow.focus();
  }
});
</script>

<template>
  <settings-menu v-model="showSettings" @cartridge:insert="cartridgeInserted($event)" />
  <menu-button v-model="showSettings" />

  <iframe
    ref="iframe"
    src="game_frame/index.html"
    class="w-full h-full overflow-hidden"
    title="Game"
  ></iframe>
</template>
