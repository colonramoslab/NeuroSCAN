const registry = new Map();

const InstancesRegistry = {
  has(uid) {
    return registry.has(uid);
  },

  get(uid) {
    return registry.get(uid);
  },

  add(simpleInstance) {
    const id = simpleInstance.wrappedObj?.id || simpleInstance.getId();
    if (!registry.has(id)) {
      registry.set(id, simpleInstance);
    }
  },

  addAll(simpleInstances) {
    simpleInstances.forEach((si) => this.add(si));
  },

  remove(uid) {
    registry.delete(uid);
  },

  getAll() {
    return Array.from(registry.values());
  },

  clear() {
    registry.clear();
  },

  size() {
    return registry.size;
  },
};

export default InstancesRegistry;
