import { v4 as uuidv4 } from 'uuid';
import SimpleInstance from '@metacell/geppetto-meta-core/model/SimpleInstance';
import {
  setOriginalColors,
  addLastSelectedInstance,
  updateWidgetConfig,
  darkenColorsFlashing,
} from '../redux/actions/widget';
import { getDataOverlay, clearDataOverlay } from '../redux/actions/dataOverlay';
import urlService from './UrlService';
import zipService from './ZipService';
import { GeppettoAdapter, InstancesRegistry } from '../infra/geppetto';
import {
  CONTACT_TYPE, CPHATE_TYPE, filesURL, NERVE_RING_TYPE, NEURON_TYPE, SYNAPSE_TYPE, SCALE_TYPE,
} from '../utilities/constants';

export const LOCAL_TYPE = 'local';

export const instanceEqualsInstance = (instanceA, instanceB) => instanceA.uid === instanceB.uid
  && instanceA.instanceType === instanceB.instanceType;

export const invertColor = ({
  r, g, b, a,
}) => ({
  r: 1 - r, g: 1 - g, b: 1 - b, a,
});

export const darkenColor = ({
  r, g, b, a,
}) => ({
  r, g, b, a: a - 0.7,
});

export const lightenColor = ({
  r, g, b, a,
}) => ({
  r, g, b, a: a + 0.7,
});

export const resetDataOverlay = (dispatch) => {
  dispatch(clearDataOverlay());
};

export const invertColorSelectedInstances = (instances, selectedUids) => (
  instances
    .map((instance) => {
      if (selectedUids.indexOf(instance.uid) < 0) {
        return instance;
      }
      const { color } = instance;
      const newInstance = { ...instance };
      if (instance.colorOriginal) {
        newInstance.color = instance.flash ? invertColor(color) : color;
      } else if (color) {
        delete newInstance.color;
      } else {
        // eslint-disable-next-line object-curly-newline
        newInstance.color = { r: 1, g: 0, b: 0, a: 1 };
      }
      return newInstance;
    }));

// darkenColorSelectedInstances is used to darken the color of the selected instances by 20%
export const darkenColorSelectedInstances = (instances, selectedUids) => (
  instances
    .map((instance) => {
      if (selectedUids.indexOf(instance.uid) < 0) {
        return instance;
      }
      const { color } = instance;
      const newInstance = { ...instance };
      newInstance.flash = !newInstance.flash;
      if (instance.colorOriginal) {
        newInstance.color = instance.flash ? darkenColor(color) : lightenColor(color);
      } else if (color) {
        delete newInstance.color;
      } else {
        // eslint-disable-next-line object-curly-newline
        newInstance.color = { r: 1, g: 0, b: 0, a: 1 };
      }
      return newInstance;
    })
);

export const setOriginalColorSelectedInstances = (instances, selectedUids) => (
  instances
    .map((instance) => {
      if (selectedUids.indexOf(instance.uid) < 0) {
        return instance;
      }
      const newInstance = {
        ...instance,
        flash: false,
        color: instance.colorOriginal ? instance.colorOriginal : instance.color,
      };
      if (!instance.colorOriginal || !instance.color) {
        delete newInstance.color;
        delete newInstance.colorOriginal;
      }
      return newInstance;
    }));
const updateInstanceSelected = (instances, selectedUids) => {
  const i = instances.map((instance) => {
    if (selectedUids.find((x) => x === instance.uid)) {
      return {
        ...instance,
        selected: true,
        flash: true,
        colorOriginal: instance.color,
        color: instance.color,
      };
    }
    return {
      ...instance,
      selected: false,
    };
  });
  return i;
};

const hideInstanceSelected = (instances, selectedUids) => instances.map((instance) => {
  if (selectedUids.find((x) => x === instance.uid)) {
    return {
      ...instance,
      hidden: true,
    };
  }
  return {
    ...instance,
  };
});

const showInstanceSelected = (instances, selectedUids) => instances.map((instance) => {
  if (selectedUids.find((x) => x === instance.uid)) {
    return {
      ...instance,
      hidden: false,
    };
  }
  return {
    ...instance,
  };
});

export const setSelectedInstances = (dispatch, viewerId, instances, selectedUids) => {
  const newInstances = updateInstanceSelected(
    instances, selectedUids,
  );

  const selected = instances.find((instance) => selectedUids.includes(instance.uid));

  if (selected && (selected.instanceType === 'scale')) {
    return;
  }

  const canShowDataOverlay = ['neuron', 'contact', 'synapse'];

  if (selected && canShowDataOverlay.includes(selected.instanceType)) {
    dispatch(getDataOverlay(selected));
  }

  const colorPickerColor = selectedUids.length > 0
    ? newInstances.find((i) => i.uid === selectedUids[selectedUids.length - 1]).colorOriginal
    : null;

  dispatch(updateWidgetConfig(
    viewerId, {
      flash: true,
      hidden: false,
      instances: newInstances,
      colorPickerColor,
    },
  ));

  let counter = 1;
  const interval = setInterval(() => {
    if (counter === 6) {
      clearInterval(interval);
      dispatch(setOriginalColors(viewerId, selectedUids));
    } else {
      dispatch(darkenColorsFlashing(viewerId, selectedUids));
    }
    counter += 1;
  }, 750);

  const newSelectedUid = instances.find((item) => (item.selected === false)
      && selectedUids.includes(item.uid));
  if (newSelectedUid) {
    dispatch(addLastSelectedInstance(viewerId, [newSelectedUid.uid]));
  }
};

export const deleteSelectedInstances = (
  dispatch, selectedInstanceToDelete, widgets, selectedUids,
) => {
  const currentWidget = widgets[selectedInstanceToDelete.viewerId];
  const { config } = currentWidget;

  const newInstances = config?.instances.filter((instance) => !selectedUids.includes(instance.uid));
  const newAddedObjectsToViewer = config?.addedObjectsToViewer
    .filter((obj) => !selectedUids.includes(obj.uid));

  dispatch(updateWidgetConfig(
    selectedInstanceToDelete.viewerId, {
      instances: newInstances,
      addedObjectsToViewer: newAddedObjectsToViewer,
    },
  ));
};

export const hideSelectedInstances = (dispatch, viewerId, instances, selectedUids) => {
  const newInstances = hideInstanceSelected(instances, selectedUids);
  dispatch(updateWidgetConfig(
    viewerId, {
      instances: newInstances,
    },
  ));
};

export const showSelectedInstances = (dispatch, viewerId, instances, selectedUids) => {
  const newInstances = showInstanceSelected(instances, selectedUids);
  dispatch(updateWidgetConfig(
    viewerId, {
      instances: newInstances,
    },
  ));
};

export const updateInstanceGroup = (instances, instanceList, newGroup = null) => instances
  .map((instance) => {
    if (instanceList.find((x) => x.uid === instance.uid)) {
      return {
        ...instance,
        group: newGroup === instance.group ? null : newGroup,
      };
    }
    return {
      ...instance,
    };
  });

export const setInstancesColor = (instances, instanceList, newColor = null) => instances
  .map((instance) => {
    if (instanceList.find((x) => x.uid === instance.uid)) {
      return {
        ...instance,
        color: newColor,
      };
    }
    return {
      ...instance,
    };
  });

const getDevStageFromTimepoint = (timepoint, devStages) => {
  const devStage = devStages
    .find((stage) => stage.timepoints !== null
      && stage.timepoints.includes(timepoint));
  return devStage.uid;
};

export const getLocationPrefixFromType = (item, devStages) => {
  if (item.instanceType === LOCAL_TYPE) {
    return 'local';
  }
  const devStage = getDevStageFromTimepoint(item.timepoint, devStages);
  switch (item.instanceType) {
    case NEURON_TYPE: {
      return `${filesURL}/neuroscan/${devStage}/${item.timepoint}/neurons/${item.filename}`;
    }
    case CONTACT_TYPE: {
      return `${filesURL}/neuroscan/${devStage}/${item.timepoint}/contacts/${item.filename}`;
    }
    case SYNAPSE_TYPE: {
      return `${filesURL}/neuroscan/${devStage}/${item.timepoint}/synapses/${item.filename}`;
    }
    case CPHATE_TYPE: {
      return `${filesURL}/neuroscan/${devStage}/${item.timepoint}/cphate/`;
    }
    case NERVE_RING_TYPE: {
      return `${filesURL}/neuroscan/${devStage}/${item.timepoint}/nerveRing/`;
    }
    case SCALE_TYPE: {
      return `${filesURL}/neuroscan/${devStage}/${item.timepoint}/scale/${item.filename}`;
    }
    default: {
      return '';
    }
  }
};

export const mapLocalGltfToInstance = ({ fileName, base64 }) => ({
  id: null,
  uid: `local_${uuidv4().replace(/-/g, '_')}`,
  uidFromDb: null,
  name: fileName.replace(/\.(gltf|glb)$/i, ''),
  selected: false,
  color: {
    r: Math.random(), g: Math.random(), b: Math.random(), a: 0.98,
  },
  instanceType: LOCAL_TYPE,
  group: null,
  content: {
    type: 'base64',
    location: 'local',
    fileName,
    base64,
  },
  getId() { return this.id; },
});

export const buildColor = (arr) => ({
  r: arr[0],
  g: arr[1],
  b: arr[2],
  a: arr[3],
});

export const mapToInstance = (item, devStages) => {
  const fileName = item.filename || '';
  const location = getLocationPrefixFromType(item, devStages);

  let color = {
    r: Math.random(), g: Math.random(), b: Math.random(), a: 0.98,
  };

  if (item.color && item.color.length === 4) {
    color = buildColor(item.color);
  }

  return {
    id: item.id,
    uid: `i_${item.uid.replace(/[-&~]/g, '_')}_${item.timepoint}`,
    uidFromDb: item.uid,
    name: item.uid,
    selected: false,
    color,
    instanceType: item.instanceType,
    group: null,
    content: {
      type: 'url',
      location,
      fileName,
    },
    getId: () => this.id,
  };
};

const getContentService = (content) => {
  switch (content.type.toLowerCase()) {
    case 'zip':
      return zipService;
    default:
      return urlService;
  }
};

const base64ToUtf8Text = (input) => {
  if (!input) return input;

  const trimmed = input.trim();
  if (trimmed.startsWith('{') || trimmed.startsWith('[')) return input;

  const b64 = trimmed.startsWith('data:')
    ? trimmed.split(',')[1]
    : trimmed;

  const binary = window.atob(b64);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i += 1) bytes[i] = binary.charCodeAt(i);

  return new TextDecoder('utf-8').decode(bytes);
};

const hasExternalReferences = (gltfJson) => {
  try {
    const parsed = JSON.parse(gltfJson);
    const externalBuffers = (parsed.buffers || []).some(
      (b) => b.uri && !b.uri.startsWith('data:'),
    );
    const externalImages = (parsed.images || []).some(
      (img) => img.uri && !img.uri.startsWith('data:'),
    );
    return externalBuffers || externalImages;
  } catch {
    return false;
  }
};

const createSimpleInstance = async (instance) => {
  const { content } = instance;
  const Resources = GeppettoAdapter.getResources();

  let base64Content;
  if (content.type && content.type.toLowerCase() === 'base64') {
    base64Content = content.base64;
  } else {
    const contentService = getContentService(content);
    base64Content = await contentService.getBase64(content.location, content.fileName);
  }
  let visualValue;
  const fileExtension = content.fileName.split('.').pop().toLowerCase();
  switch (fileExtension) {
    case 'obj':
      visualValue = {
        eClass: Resources.OBJ,
        obj: base64Content,
      };
      break;
    case 'gltf': {
      const gltfJsonText = base64ToUtf8Text(base64Content);
      if (hasExternalReferences(gltfJsonText)) {
        console.warn(`GLTF file "${content.fileName}" references external .bin or texture files which cannot be resolved from a local drag-and-drop.`);
      }
      visualValue = {
        eClass: Resources.GLTF,
        gltf: gltfJsonText,
      };
      break;
    }
    case 'glb': {
      visualValue = {
        eClass: Resources.GLTF,
        gltf: base64Content,
      };
      break;
    }
    default:
      visualValue = {
        eClass: Resources.OBJ,
        obj: base64Content,
      };
  }

  return new SimpleInstance({
    eClass: 'SimpleInstance',
    id: instance.uid,
    name: instance.uid,
    type: { eClass: 'SimpleType' },
    visualValue,
  });
};

// this project runs on node 14, so settimeout is not wrapped in a promise yet.
const delay = (ms) => new Promise((res) => setTimeout(res, ms));

export const createSimpleInstancesFromInstances = (instances) => {
  const newInstances = instances.filter(
    (instance) => !InstancesRegistry.has(instance.uid),
  );
  if (newInstances.length > 100) {
    const results = [];
    const batches = [];
    const batchSize = 100;
    for (let i = 0; i < newInstances.length; i += batchSize) {
      batches.push(newInstances.slice(i, i + batchSize));
    }

    // eslint-disable-next-line no-async-promise-executor
    return new Promise(async (resolve) => {
      while (batches.length >= 1) {
        // eslint-disable-next-line no-loop-func, no-await-in-loop
        results.push(await Promise.all(
          batches[0].map((instance) => createSimpleInstance(instance)),
        ));
        // eslint-disable-next-line no-await-in-loop
        await delay(5000);
        batches.shift();
      }
      resolve(results);
    }).then((newSimpleInstances) => {
      GeppettoAdapter.addInstances(newSimpleInstances.flat());
    });
  }

  return Promise.all(
    newInstances.map((instance) => createSimpleInstance(instance)),
  ).then((newSimpleInstances) => {
    GeppettoAdapter.addInstances(newSimpleInstances);
  });
};

export const getGroupsFromInstances = (instances) => (
  [
    ...new Set(
      instances
        .filter((instance) => instance.group)
        .map((instance) => instance.group),
    ),
  ]);

export const groupBy = (items, key) => {
  const results = {};
  for (let i = 0; items.length > i; i += 1) {
    if (items[i][key] === null) {
      // eslint-disable-next-line no-continue
      continue;
    }
    if (!(items[i][key] in results)) {
      results[items[i][key]] = [];
    }
    results[items[i][key]].push(items[i]);
  }
  return results;
};

export const getInstancesOfType = (instances, instanceType) => (
  groupBy(instances, 'instanceType'))[instanceType];

export const getInstancesByGroups = (instances) => (
  groupBy(instances, 'group'));

export const handleSelect = (dispatch, viewerId, selectedInstance, widgets) => {
  if (viewerId) {
    const { instances } = widgets[viewerId].config;
    setSelectedInstances(dispatch, viewerId, instances, [selectedInstance.uid]);
  }
};

export const sortedInstances = (instances) => {
  instances.sort((a, b) => {
    const neuronsA = a.neurons.join('');
    const neuronsB = b.neurons.join('');
    return neuronsA.localeCompare(neuronsB);
  });

  instances.forEach((obj) => {
    if (obj.neurons.length > 1) {
      obj.neurons.sort();
    }
  });
  return instances;
};

export const sortedGroupedIterations = (items) => items.map((innerArray) => innerArray.slice()
  .sort((a, b) => a.name.localeCompare(b.name)));

export async function fetchDataForEntity(service, timePoint, entities) {
  if (entities.length === 0) return [];
  const newEntities = await service.getByUID(timePoint, entities.map((e) => e.uidFromDb));
  return newEntities;
}
