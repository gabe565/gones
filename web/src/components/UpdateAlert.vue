<template>
  <transition name="slide-fade" appear>
    <section
      role="alert"
      v-if="offlineReady || needRefresh"
      class="z-40 flex items-center absolute bottom-0 right-0 m-8 mb-md-2 bg-teal-700 border-b-4 border-teal-500 rounded-lg text-teal-100 px-4 py-3 shadow-md"
    >
      <span class="mx-1">
        <template v-if="needRefresh"
          >New content available, click on reload button to update.</template
        >
        <template v-else>App ready to work offline</template>
      </span>
      <button
        v-if="needRefresh"
        class="mx-1 p-2 bg-blue-800 hover:bg-blue-700 rounded-full transition-colors"
        :class="{ 'animate-spin': loading }"
        @click="updateServiceWorker"
      >
        <icon-refresh-cw />
        <span class="sr-only">Refresh</span>
      </button>
      <button
        class="mx-1 p-2 bg-teal-600 hover:bg-teal-500 rounded-full transition-colors"
        @click="closePromptUpdateSW"
      >
        <icon-close />
        <span class="sr-only">Close</span>
      </button>
    </section>
  </transition>
</template>

<script setup>
import IconRefreshCw from "~icons/feather/refresh-cw";
import IconClose from "~icons/feather/x";
import { onMounted, ref, watch } from "vue";

const updateSW = ref(undefined);
const offlineReady = ref(false);
const needRefresh = ref(false);
const loading = ref(false);

const intervalMS = 60 * 60 * 1000;

onMounted(async () => {
  try {
    const { registerSW } = await import("virtual:pwa-register");
    updateSW.value = registerSW({
      immediate: true,
      onOfflineReady() {
        offlineReady.value = true;
        console.log("onOfflineReady");
      },
      onNeedRefresh() {
        needRefresh.value = true;
        console.log("onNeedRefresh");
      },
      onRegistered(swRegistration) {
        if (swRegistration) {
          setInterval(() => swRegistration.update(), intervalMS);
        }
      },
      onRegisterError(error) {
        console.error(error);
      },
    });
  } catch {
    console.log("PWA disabled.");
  }
});

const closePromptUpdateSW = async () => {
  offlineReady.value = false;
  needRefresh.value = false;
};

const updateServiceWorker = () => {
  if (updateSW.value) {
    loading.value = true;
    updateSW.value(true);
  }
};

watch(offlineReady, (val) => {
  if (val && !needRefresh.value) {
    setTimeout(() => (offlineReady.value = false), 5000);
  }
});
</script>

<style scoped lang="scss">
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 400ms cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateY(0.5em);
  opacity: 0;
}
</style>
