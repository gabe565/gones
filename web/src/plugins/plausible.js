import Plausible from "plausible-tracker";

export const plausible = Plausible({
  apiHost: import.meta.env.VITE_PLAUSIBLE_HOST,
});

if (import.meta.env.VITE_PLAUSIBLE_HOST) {
  plausible.enableAutoPageviews();
  plausible.enableAutoOutboundTracking();
} else {
  plausible.trackPageview = () => {};
  plausible.trackEvent = () => {};
}
