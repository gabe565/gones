<script setup>
import IconClose from "~icons/material-symbols/close-rounded";
import IconUpload from "~icons/material-symbols/upload-rounded";
import IconGithub from "~icons/simple-icons/github";
import KeyTable from "./KeyTable.vue";
import { ref } from "vue";

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
      class="flex flex-col fixed h-full w-[22em] top-0 left-0 z-30 bg-gray-900 p-5 shadow-lg shadow-gray-950"
      role="menu"
      aria-expanded="true"
    >
      <div class="flex items-start justify-between pb-6">
        <h1 class="text-4xl pb-2">GoNES</h1>

        <button
          class="p-1 bg-gray-800 hover:bg-gray-700 rounded-full transition-colors"
          @click.prevent="$emit('update:modelValue', !modelValue)"
        >
          <icon-close aria-hidden="true" />
          <span class="sr-only">Close menu</span>
        </button>
      </div>

      <div class="flex flex-col items-center pb-6">
        <input
          ref="cartridgeInput"
          type="file"
          class="hidden"
          @change="$emit('cartridge:insert', $event.target.files.item(0))"
        />
        <button
          class="block bg-gray-800 hover:bg-gray-700 py-2 px-4 my-2 rounded-full border border-gray-600 transition-colors"
          @click="cartridgeInput.click()"
        >
          <icon-upload class="inline -mt-0.5" aria-hidden="true" />
          Load Game
        </button>
      </div>

      <div class="flex-grow" />

      <div class="w-11/12 self-center py-2.5">
        <key-table />
      </div>

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
