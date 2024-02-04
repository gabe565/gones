<script setup>
defineProps({
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
});
</script>

<template>
  <button
    class="block bg-gray-800 hover:bg-gray-700 rounded-full border border-gray-700 transition-colors"
    :class="[icon ? 'p-1' : 'py-2 px-4']"
  >
    <slot name="prepend">
      <component
        :is="prependIcon"
        v-if="prependIcon"
        class="inline -mt-0.5 mr-1.5"
        aria-hidden="true"
      />
    </slot>
    <slot name="icon">
      <template v-if="icon">
        <component :is="icon" v-if="icon" />
      </template>
    </slot>
    <slot>
      <span :class="{ 'sr-only': icon }">{{ text }}</span>
    </slot>
  </button>
</template>
