<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import MenuButton from "./MenuButton.vue";
import SettingsMenu from "./SettingsMenu.vue";
import {
  exitEvent,
  nameEvent,
  newEventPromise,
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
let { promise: ready, resolve: readyResolve } = newEventPromise();
const handleReady = () => readyResolve();

// Promise that will resolve when the game exits
let { promise: exit, resolve: exitResolve } = newEventPromise();
const handleExit = () => exitResolve();

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
  window.addEventListener(exitEvent, handleExit);
});

onBeforeUnmount(() => {
  window.removeEventListener(readyEvent, handleReady);
  window.removeEventListener(nameEvent, handleName);
  window.removeEventListener(exitEvent, handleExit);
});

const iframe = ref();
let running = ref(false);

const cartridgeInserted = async (val) => {
  showSettings.value = false;
  if (running.value) {
    iframe.value.contentWindow.dispatchEvent(newExitEvent());
    await exit;
    ({ promise: ready, resolve: readyResolve } = newEventPromise());
    ({ promise: exit, resolve: exitResolve } = newEventPromise());
    await iframe.value.contentWindow.location.reload();
    await ready;
  }
  iframe.value.contentWindow.dispatchEvent(newPlayEvent(val.name, val.arrayBuffer()));
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
  await exit;
  ({ promise: exit, resolve: exitResolve } = newEventPromise());
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

  <iframe ref="iframe" src="game.html" class="w-full h-full overflow-hidden" title="Game" />
</template>
