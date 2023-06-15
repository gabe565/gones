<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from "vue";

const props = defineProps({
  cartridge: {
    type: File,
    default: null,
  },
  showSettings: {
    type: Boolean,
  },
});

// Promise that will resolve when the iframe is done reloading
let resolve;
let promise = new Promise((r) => {
  resolve = r;
});

// Handle messages from the iframe
const iframeMessage = ({ data }) => {
  if (data.type === "ready") {
    resolve();
  }
};
onMounted(() => {
  window.addEventListener("message", iframeMessage);
});
onBeforeUnmount(() => {
  window.removeEventListener("message", iframeMessage);
});

const iframe = ref();
let running = false;

watch(
  () => props.cartridge,
  async (val) => {
    if (running) {
      running = true;
      promise = new Promise((r) => {
        resolve = r;
      });
      iframe.value.contentWindow.location.reload();
      await promise;
    }
    iframe.value.contentWindow.postMessage({
      type: "play",
      cartridge: new Uint8Array(await val.arrayBuffer()),
    });
    iframe.value.contentWindow.focus();
    running = true;
  },
);

watch(
  () => props.showSettings,
  (val) => {
    if (!val && running) {
      iframe.value.contentWindow.focus();
    }
  },
);
</script>

<template>
  <iframe
    ref="iframe"
    src="game_frame/index.html"
    class="w-full h-full overflow-hidden"
    title="Game"
  ></iframe>
</template>
