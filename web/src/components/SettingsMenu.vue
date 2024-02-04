<script setup>
import IconClose from "~icons/material-symbols/close-rounded";
import IconOpen from "~icons/material-symbols/folder-open-rounded";
import IconGithub from "~icons/simple-icons/github";
import IconLogo from "~icons/gones/icon";
import IconHeading from "~icons/gones/heading?width=8em&height=2.5em";
import KeyTable from "./KeyTable.vue";
import { ref } from "vue";
import GonesButton from "./GonesButton.vue";

defineProps({
  modelValue: {
    type: Boolean,
    default: true,
  },
});

defineEmits(["update:modelValue", "cartridge:insert"]);
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
          <icon-heading alt="GoNES" />
        </div>

        <div class="flex-grow" />

        <gones-button
          text="Close menu"
          :icon="IconClose"
          @click.prevent="$emit('update:modelValue', !modelValue)"
        />
      </div>

      <div class="flex flex-col items-center pb-6">
        <input
          ref="cartridgeInput"
          type="file"
          class="hidden"
          accept=".nes"
          @change="$emit('cartridge:insert', $event.target.files.item(0))"
        />
        <gones-button text="Open ROM" :prepend-icon="IconOpen" @click="cartridgeInput.click()" />
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
