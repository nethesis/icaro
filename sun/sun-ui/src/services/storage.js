var StorageService = {
  methods: {
    get(prop) {
      return JSON.parse(localStorage.getItem(prop));
    },
    set(prop, object) {
      localStorage.setItem(prop, JSON.stringify(object));
    },
    delete(prop) {
      localStorage.removeItem(prop);
    },
    clear() {
      localStorage.clear();
    }
  }
};
export default StorageService;
