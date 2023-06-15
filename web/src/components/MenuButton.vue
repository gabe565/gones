<script setup>
import IconMenu from "~icons/feather/menu";

defineProps({
  modelValue: {
    type: Boolean,
    default: true,
  },
});

defineEmits(["update:modelValue"]);
</script>

<template>
  <button
    class="text-xl p-3 absolute top-0 left-0 z-10 opacity-25 hover:opacity-70 transition-opacity"
    @click.prevent="$emit('update:modelValue', !modelValue)"
  >
    <icon-menu aria-hidden="true" />
    <span class="sr-only">{{ modelValue ? "Hide" : "Show" }} menu</span>
  </button>
</template>
