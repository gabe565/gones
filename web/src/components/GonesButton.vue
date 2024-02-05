<script setup>
import { computed } from "vue";

const props = defineProps({
  text: {
    type: String,
    default: "",
  },
  prependIcon: {
    type: Object,
    default: null,
  },
  icon: {
    type: Object,
    default: null,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  size: {
    type: String,
    default: "medium",
  },
});

const padClass = computed(() => {
  if (props.icon) {
    return "p-1";
  }
  switch (props.size) {
    case "x-small":
      return "px-2.5";
    case "small":
      return "py-1 px-3";
  }
  return "py-2 px-4";
});

const textClass = computed(() => {
  switch (props.size) {
    case "x-small":
      return "text-xs";
    case "small":
      return "text-sm";
  }
  return "text-base";
});

const iconClass = computed(() => {
  switch (props.size) {
    case "x-small":
    case "small":
      return "mr-1";
  }
  return "mr-1.5";
});
</script>

<template>
  <button
    class="block rounded-full border border-gray-700 transition-colors"
    :class="[
      disabled ? 'text-gray-500 bg-gray-700' : 'bg-gray-800 hover:bg-gray-700',
      padClass,
      textClass,
    ]"
    :disabled="disabled"
  >
    <slot name="prepend">
      <component
        :is="prependIcon"
        v-if="prependIcon"
        class="inline -mt-0.5"
        :class="[iconClass]"
        aria-hidden="true"
      />
    </slot>
    <slot name="icon">
      <template v-if="icon">
        <component :is="icon" v-if="icon" />
      </template>
    </slot>
    <slot>
      <span :class="[{ 'sr-only': icon }]">{{ text }}</span>
    </slot>
  </button>
</template>
