export const newEventPromise = () => {
  let resolve;
  const promise = new Promise((r) => (resolve = r));
  return { promise, resolve };
};

export const readyEvent = "gonesReady";
export const newReadyEvent = () => new CustomEvent(readyEvent);

export const playEvent = "gonesPlay";
export const newPlayEvent = (cartridge) => new CustomEvent(playEvent, { detail: { cartridge } });

export const nameEvent = "gonesName";
export const newNameEvent = (value) => new CustomEvent(nameEvent, { detail: { value } });

export const exitEvent = "gonesExit";
export const newExitEvent = () => new CustomEvent(exitEvent);

export const saveStateEvent = "gonesSaveState";
export const newSaveStateEvent = () => new CustomEvent(saveStateEvent);

export const loadStateEvent = "gonesloadState";
export const newLoadStateEvent = () => new CustomEvent(loadStateEvent);
