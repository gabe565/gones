<script setup>
import { ref } from "vue";

const keys = ref({
  "A button": ["M"],
  "B button": ["N"],
  "Start button": ["Enter"],
  "Select button": ["RightShift"],
  "Arrow buttons": ["W", "A", "S", "D"],
  "A turbo": ["K"],
  "B turbo": ["J"],
  "Save State": ["F1"],
  "Load State": ["F5"],
});
</script>

<template>
  <table class="w-full mb-1">
    <caption class="text-2xl pb-1">
      Keybinds
    </caption>
    <thead class="border-b border-gray-700">
      <tr>
        <th>Key</th>
        <th class="text-left">Action</th>
      </tr>
    </thead>
    <tbody>
      <tr
        v-for="(keyList, action) in keys"
        :key="action"
        class="hover:bg-gray-800 transition-colors"
      >
        <td class="text-center">
          <kbd
            v-for="key in keyList"
            :key="key"
            class="px-1.5 py-0.5 m-0.5 text-xs font-semibold bg-gray-700 text-gray-100 border-gray-500 border border-b-2 rounded-lg"
            >{{ key }}</kbd
          >
        </td>
        <td class="text-sm">{{ action }}</td>
      </tr>
    </tbody>
  </table>
  <hr class="border-gray-700" />
</template>
