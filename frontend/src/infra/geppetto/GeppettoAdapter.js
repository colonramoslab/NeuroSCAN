import InstancesRegistry from './InstancesRegistry';

let geppettoResources = null;
let geppettoManager = null;
let initialized = false;

const GeppettoAdapter = {
  init() {
    if (initialized) return;

    // eslint-disable-next-line global-require
    const Manager = require('@metacell/geppetto-meta-client/common/Manager').default;
    // eslint-disable-next-line global-require
    const Resources = require('@metacell/geppetto-meta-core/Resources').default;

    const GEPPETTO = {};
    GEPPETTO.Resources = Resources;
    GEPPETTO.Manager = new Manager();

    window.GEPPETTO = GEPPETTO;
    window.Instances = [];

    geppettoResources = Resources;
    geppettoManager = GEPPETTO.Manager;
    initialized = true;
  },

  getResources() {
    return geppettoResources;
  },

  getManager() {
    return geppettoManager;
  },

  syncInstancesToWindow() {
    window.Instances = InstancesRegistry.getAll();
    if (geppettoManager) {
      geppettoManager.augmentInstancesArray(window.Instances);
    }
  },

  addInstances(simpleInstances) {
    InstancesRegistry.addAll(simpleInstances);
    this.syncInstancesToWindow();
  },

  hasInstance(uid) {
    return InstancesRegistry.has(uid);
  },
};

export default GeppettoAdapter;
