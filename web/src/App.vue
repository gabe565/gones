<script setup>
import { ref } from "vue";
import SettingsMenu from "./components/SettingsMenu.vue";
import MenuButton from "./components/MenuButton.vue";
import GameView from "./components/GameView.vue";
import UpdateAlert from "./components/UpdateAlert.vue";

const showSettings = ref(true);
const cartridge = ref(null);

const cartridgeInserted = (val) => {
  cartridge.value = val;
  showSettings.value = false;
};
</script>

<template>
  <update-alert />
  <settings-menu v-model="showSettings" @cartridge:insert="cartridgeInserted($event)" />
  <menu-button v-model="showSettings" />
  <game-view :cartridge="cartridge" :show-settings="showSettings" />
</template>
