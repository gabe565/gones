<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import MenuButton from "./MenuButton.vue";
import SettingsMenu from "./SettingsMenu.vue";
import { wait } from "../util/wait";

const showSettings = ref(true);
const name = ref("");
const defaultTitle = document.title;

// Promise that will resolve when the iframe is done reloading
let resolve;
let promise = new Promise((r) => {
  resolve = r;
});

// Handle messages from the iframe
const iframeMessage = ({ data }) => {
  if (data.type === "ready") {
    resolve();
  } else if (data.type === "name") {
    name.value = data.value;
    if (data.value) {
      document.title = data.value + " - " + defaultTitle;
    } else {
      document.title = defaultTitle;
    }
  }
};
onMounted(() => {
  window.addEventListener("message", iframeMessage);
});
onBeforeUnmount(() => {
  window.removeEventListener("message", iframeMessage);
});

const iframe = ref();
let running = ref(false);

const cartridgeInserted = async (val) => {
  showSettings.value = false;
  if (running.value) {
    running.value = true;
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
  running.value = true;
};

watch(showSettings, (val) => {
  if (!val && running) {
    iframe.value.contentWindow.focus();
  }
});

const stop = async () => {
  running.value = false;
  iframe.value.contentWindow.Gones.exit();
  await wait(100);
  await iframe.value.contentWindow.location.reload();
  name.value = "";
  document.title = defaultTitle;
};

const saveState = () => {
  iframe.value.contentWindow.Gones.saveState();
  showSettings.value = false;
};

const loadState = () => {
  iframe.value.contentWindow.Gones.loadState();
  showSettings.value = false;
};
</script>

<template>
  <settings-menu
    v-model="showSettings"
    :running="running"
    :name="name"
    @gones:cartridge="cartridgeInserted($event)"
    @gones:stop="stop"
    @gones:save-state="saveState"
    @gones:load-state="loadState"
  />
  <menu-button v-model="showSettings" />

  <iframe
    ref="iframe"
    src="game_frame/index.html"
    class="w-full h-full overflow-hidden"
    title="Game"
  ></iframe>
</template>
