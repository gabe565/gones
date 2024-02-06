export const DbName = "gones";
export const DbVersion = 1;

export const getDb = async () => {
  return new Promise((resolve, reject) => {
    let request = window.indexedDB.open(DbName, DbVersion);

    request.onerror = reject;
    request.onsuccess = ({ target: { result } }) => {
      resolve(result);
    };

    request.onupgradeneeded = ({ target: { result: db } }) => {
      console.log("Upgrading", DbName, "db to v" + DbVersion);

      db.createObjectStore("states", { keyPath: "name" });
      const store = db.createObjectStore("saves", { keyPath: "name" });

      store.transaction.oncomplete = () => {
        const stateStore = db.transaction("states", "readwrite").objectStore("states");
        const saveStore = db.transaction("saves", "readwrite").objectStore("saves");
        for (let i = 0; i < localStorage.length; i++) {
          const name = localStorage.key(i);
          const data = localStorage.getItem(name);
          if (name.endsWith(".state.gz")) {
            console.log("Migrating state", name);
            stateStore.add({ name, data });
          } else if (name.endsWith(".sav")) {
            console.log("Migrating save", name);
            saveStore.add({ name, data });
          }
        }
      };
    };
  });
};

export const dbPut = async (table, name, data) => {
  const db = await getDb();

  return new Promise((resolve) => {
    const store = db.transaction([table], "readwrite").objectStore(table);
    const request = store.put({ name, data });
    request.onsuccess = resolve;
  });
};

export const dbGet = async (table, name) => {
  const db = await getDb();

  return new Promise((resolve) => {
    const store = db.transaction(table).objectStore(table);
    store.get(name).onsuccess = ({ target }) => {
      if (target.result) {
        resolve(target.result.data);
      } else {
        resolve(null);
      }
    };
  });
};
