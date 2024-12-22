<script setup>
import { ref } from "vue";
import GonesButton from "./GonesButton.vue";
import KeyTable from "./KeyTable.vue";
import IconHeading from "~icons/gones/heading?width=8em&height=2.5em";
import IconLogo from "~icons/gones/icon";
import IconClose from "~icons/material-symbols/close-rounded";
import IconOpen from "~icons/material-symbols/folder-open-rounded";
import IconLoad from "~icons/material-symbols/restore-page-rounded";
import IconSave from "~icons/material-symbols/save-rounded";
import IconStop from "~icons/material-symbols/stop-rounded";
import IconGithub from "~icons/simple-icons/github";

defineProps({
  modelValue: {
    type: Boolean,
    default: true,
  },
  running: {
    type: Boolean,
    default: false,
  },
  name: {
    type: String,
    default: "",
  },
});

defineEmits([
  "update:modelValue",
  "gones:cartridge",
  "gones:stop",
  "gones:saveState",
  "gones:loadState",
]);
const cartridgeInput = ref();
</script>

<template>
  <transition name="slide-fade">
    <section
      v-if="modelValue"
      id="settings"
      class="flex flex-col fixed h-full w-[22em] top-0 left-0 z-30 bg-gray-900 pb-5 shadow-lg shadow-gray-950"
      role="menu"
      aria-expanded="true"
    >
      <div
        class="h-16 flex items-center px-3 mb-7 bg-gradient-to-b from-gray-900 to-gray-950 border-b border-gray-700"
      >
        <div class="h-full p-3 pl-0 border-r border-gray-700">
          <icon-logo class="text-4xl" aria-hidden="true" />
        </div>
        <div class="p-3">
          <icon-heading aria-hidden="true" />
          <h1 class="sr-only">GoNES</h1>
        </div>

        <div class="flex-grow" />

        <gones-button
          text="Close menu"
          :icon="IconClose"
          @click.prevent="$emit('update:modelValue', !modelValue)"
        />
      </div>

      <div class="flex flex-col items-center pb-6 gap-2 text-center">
        <h2>Game</h2>
        <div v-if="name" class="w-full px-3 truncate" :title="name"><b>Name:</b> {{ name }}</div>
        <div class="flex justify-center mb-3">
          <input
            ref="cartridgeInput"
            type="file"
            class="hidden"
            accept=".nes"
            @change="$emit('gones:cartridge', $event.target.files.item(0))"
          />
          <gones-button
            text="Open ROM"
            :prepend-icon="IconOpen"
            class="rounded-r-none border-r-0"
            @click="cartridgeInput.click()"
          />
          <gones-button
            :disabled="!running"
            text="Stop"
            :prepend-icon="IconStop"
            class="rounded-l-none"
            @click="$emit('gones:stop')"
          />
        </div>

        <h2>State</h2>
        <div class="flex justify-center">
          <gones-button
            :disabled="!running"
            :prepend-icon="IconSave"
            text="Save State"
            size="small"
            class="rounded-r-none border-r-0"
            @click="$emit('gones:saveState')"
          />
          <gones-button
            :disabled="!running"
            :prepend-icon="IconLoad"
            text="Load State"
            size="small"
            class="rounded-l-none"
            @click="$emit('gones:loadState')"
          />
        </div>
      </div>

      <div class="flex-grow" />

      <key-table />

      <p class="text-center text-gray-300 text-sm mt-2">
        <a href="https://github.com/gabe565/gones" target="_blank">
          <icon-github class="inline -mt-0.5" aria-hidden="true" />
          View on GitHub
        </a>
      </p>
    </section>
  </transition>
  <transition name="fade">
    <div
      v-if="modelValue"
      class="absolute top-0 left-0 w-full h-full bg-black opacity-70"
      @click="$emit('update:modelValue', false)"
    />
  </transition>
</template>

<style scoped lang="scss">
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 250ms cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(-100%);
  opacity: 50%;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 250ms cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
