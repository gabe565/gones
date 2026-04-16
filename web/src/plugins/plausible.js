import { init, track } from "@plausible-analytics/tracker";

const apiHost = import.meta.env.VITE_PLAUSIBLE_HOST;

if (apiHost) {
  init({
    domain: window.location.hostname,
    endpoint: apiHost + "/api/event",
    outboundLinks: true,
  });
  track("pageview");
}

export const plausible = {
  trackEvent: (eventName, options) => {
    if (apiHost) {
      track(eventName, options);
    }
  },
  trackPageview: (options) => {
    if (apiHost) {
      track("pageview", options);
    }
  },
};
