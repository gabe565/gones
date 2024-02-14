<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import MenuButton from "./MenuButton.vue";
import SettingsMenu from "./SettingsMenu.vue";
import { wait } from "../util/wait";
import {
  nameEvent,
  newExitEvent,
  newLoadStateEvent,
  newPlayEvent,
  newSaveStateEvent,
  readyEvent,
} from "../util/events";

const showSettings = ref(true);
const name = ref("");
const defaultTitle = document.title;

// Promise that will resolve when the iframe is done reloading
let resolve;
let promise = new Promise((r) => {
  resolve = r;
});

const handleReady = () => resolve();

const handleName = ({ detail: { value } }) => {
  name.value = value;
  if (value) {
    document.title = value + " - " + defaultTitle;
  } else {
    document.title = defaultTitle;
  }
};

onMounted(() => {
  window.addEventListener(readyEvent, handleReady);
  window.addEventListener(nameEvent, handleName);
});

onBeforeUnmount(() => {
  window.removeEventListener(readyEvent, handleReady);
  window.removeEventListener(nameEvent, handleName);
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
    iframe.value.contentWindow.dispatchEvent(newExitEvent());
    await wait(100);
    await iframe.value.contentWindow.location.reload();
    await promise;
  }
  iframe.value.contentWindow.dispatchEvent(newPlayEvent(val.arrayBuffer()));
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
  iframe.value.contentWindow.dispatchEvent(newExitEvent());
  await wait(100);
  await iframe.value.contentWindow.location.reload();
  name.value = "";
  document.title = defaultTitle;
};

const saveState = () => {
  iframe.value.contentWindow.dispatchEvent(newSaveStateEvent());
  showSettings.value = false;
};

const loadState = () => {
  iframe.value.contentWindow.dispatchEvent(newLoadStateEvent());
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
