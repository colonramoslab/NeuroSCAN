import * as layoutActions from '@metacell/geppetto-meta-client/common/layout/actions';
import { updateWidget } from '@metacell/geppetto-meta-client/common/layout/actions';
import { ADD_DEVSTAGES, receivedDevStages } from './actions/devStages';
import { GET_DATA_OVERLAY, renderDataOverlay } from './actions/dataOverlay';
import {
  loading,
  raiseError,
  loadingSuccess,
} from './actions/misc';
import {
  ADD_INSTANCES,
  ADD_INSTANCES_TO_GROUP,
  CLONE_VIEWER_WITH_INSTANCES_LIST,
  SET_INSTANCES_COLOR,
  UPDATE_TIMEPOINT_VIEWER,
  UPDATE_BACKGROUND_COLOR_VIEWER,
  UPDATE_WIDGET_CONFIG,
  ROTATE_START_ALL,
  ROTATE_STOP_ALL,
  updateWidgetConfig,
  INVERT_COLORS_FLASHING,
  DARKEN_COLORS_FLASHING,
  SET_ORIGINAL_COLORS_FLASHING,
  TOGGLE_INSTANCE_HIGHLIGHT,
  DELETE_FROM_WIDGET,
} from './actions/widget';
import { DevStageService } from '../services/DevStageService';
import neuronService from '../services/NeuronService';
import contactService from '../services/ContactService';
import synapseService from '../services/SynapseService';
import cphateService from '../services/CphateService';
import scaleService from '../services/ScaleService';
import nerveRingService from '../services/NerveRingService';
import {
  CONTACT_TYPE,
  NEURON_TYPE,
  SYNAPSE_TYPE,
  SCALE_TYPE,
  NERVE_RING_TYPE,
  VIEWERS,
} from '../utilities/constants';
import { cameraControlsRotateState } from '../components/Chart/CameraControls';
import { addToWidget } from '../utilities/functions';
import {
  createSimpleInstancesFromInstances,
  updateInstanceGroup,
  setInstancesColor,
  getInstancesOfType,
  mapToInstance,
  invertColorSelectedInstances,
  setOriginalColorSelectedInstances, fetchDataForEntity, darkenColorSelectedInstances,
} from '../services/instanceHelpers';

const devStagesService = new DevStageService();

export const createWidget = (store, timePoint, viewerType) => {
  const state = store.getState();
  const { widgets } = state;
  const devStages = state.devStages.neuroSCAN;
  const devStage = devStages.find((ds) => ds.begin <= timePoint && ds.end >= timePoint);
  const viewerNumber = Object.values(widgets).reduce((maxViewerNumber, w) => {
    const found = w.name.match('^Viewer (?<id>\\d+) .*');
    if (found && found.length > 0) {
      const thisViewerNumber = parseInt(found[1], 10);
      return Math.max(thisViewerNumber + 1, maxViewerNumber);
    }
    return maxViewerNumber;
  }, 1);
  const devStageLabel = devStage ? `${devStage.uid} ${timePoint}` : 'Local';
  const viewerName = `${viewerType} ${viewerNumber} (${devStageLabel})`;
  return {
    id: null,
    name: viewerName,
    type: viewerType,
    timePoint,
  };
};

const getWidget = (store, viewerId) => {
  const state = store.getState();
  const { widgets } = state;
  if (!widgets) {
    return false;
  }
  const widget = widgets[viewerId];
  return !widget ? false : { ...widget };
};

const findWidget = (store, timepoint, viewerType) => {
  const state = store.getState();
  const { widgets } = state;
  if (!widgets) {
    return false;
  }
  const widget = Object
    .values(widgets)
    .find((w) => w.timePoint === timepoint && w.type === viewerType);
  return !widget ? false : { ...widget };
};

const cleanEmpty = (store) => {
  const state = store.getState();
  const { widgets } = state;
  if (!widgets) {
    return false;
  }
  const cleanedWidgets = Object
    .values(widgets)
    .filter((w) => w.config?.instances?.length > 0);
  cleanedWidgets.forEach((w) => {
    const { addedObjectsToViewer } = w.config;
    store.dispatch(addToWidget(w, [], true, addedObjectsToViewer));
  });
  return cleanedWidgets.length > 0 ? cleanedWidgets : false;
};

const middleware = (store) => (next) => async (action) => {
  switch (action.type) {
    case ADD_DEVSTAGES: {
      const msg = 'Getting development stages';
      next(loading(msg, action.type));
      devStagesService.getDevStages().then((stages) => {
        store.dispatch(receivedDevStages(stages));
        next(loadingSuccess(msg, action.type));
      }, (error) => {
        console.error(msg, error);
        next(raiseError(msg));
      });
      break;
    }

    case ADD_INSTANCES: {
      const msg = 'Creating and adding instances to the viewer';
      next(loading(msg, action.type));
      createSimpleInstancesFromInstances(action.instances)
        .then(async () => {
          const state = store.getState();
          const timepoint = state.search.filters.timePoint;
          let widget = getWidget(store, action.viewerId, action.viewerType);
          if (!widget) {
            widget = createWidget(store, timepoint, action.viewerType);
          }
          const filteredNewInstances = Array.isArray(widget?.config?.instances)
          && widget?.config?.instances.length !== 0
            ? action.instances.filter((item2) => !widget.config.instances
              .some((item1) => item1.uid === item2.uid)) : action.instances;

          const addedObjectsToViewer = Array.isArray(widget?.config?.instances)
          && widget?.config?.instances.length !== 0
            ? widget.config.instances.concat(filteredNewInstances) : action.instances;
          store.dispatch(
            addToWidget(
              widget,
              filteredNewInstances,
              false,
              addedObjectsToViewer,
            ),
          );
          next(loadingSuccess(msg, action.type));
        }, (error) => {
          console.error(msg, error);
          next(raiseError(msg));
        });
      break;
    }

    case CLONE_VIEWER_WITH_INSTANCES_LIST: {
      const msg = 'Cloning viewer and adding instances to the viewer';
      next(loading(msg, action.type));
      const currentWidget = getWidget(store, action.fromViewerId);
      const { timePoint } = currentWidget.config;

      createSimpleInstancesFromInstances(action.instances)
        .then(() => {
          const widget = createWidget(store, timePoint, VIEWERS.InstanceViewer);
          store.dispatch(
            addToWidget(
              widget,
              action.instances,
              false,
              action.instances,
            ),
          );
          next(loadingSuccess(msg, action.type));
        }, (error) => {
          console.error(msg, error);
          next(raiseError(msg));
        });
      break;
    }

    case UPDATE_BACKGROUND_COLOR_VIEWER: {
      const widget = getWidget(store, action.viewerId);
      widget.config.backgroundColor = action.backgroundColor;
      store.dispatch(updateWidget(widget));
      break;
    }

    case UPDATE_WIDGET_CONFIG: {
      const widget = getWidget(store, action.viewerId);
      widget.config = {
        ...widget.config,
        ...action.config,
      };
      store.dispatch(updateWidget(widget));
      break;
    }

    case INVERT_COLORS_FLASHING: {
      const widget = getWidget(store, action.viewerId);
      const instances = invertColorSelectedInstances(
        widget.config.instances,
        action.uids,
      );
      widget.config = {
        ...widget.config,
        instances,
      };
      store.dispatch(updateWidget(widget));
      break;
    }

    case DARKEN_COLORS_FLASHING: {
      const widget = getWidget(store, action.viewerId);
      const instances = darkenColorSelectedInstances(
        widget.config.instances,
        action.uids,
      );
      widget.config = {
        ...widget.config,
        instances,
      };
      store.dispatch(updateWidget(widget));
      break;
    }

    case GET_DATA_OVERLAY: {
      if (!action?.instance?.instanceType || !action?.instance?.id) {
        break;
      }

      const id = action?.instance?.id;
      const msg = 'Getting data overlay';
      next(loading(msg, action.type));

      switch (action.instance?.instanceType) {
        case 'neuron':
          neuronService.getByID(id).then((neuron) => {
            store.dispatch(renderDataOverlay(neuron));
            next(loadingSuccess(msg, action.type));
          }, (error) => {
            console.error(msg, error);
            next(raiseError(error));
          });
          break;
        case 'synapse':
          synapseService.getByID(id).then((synapse) => {
            store.dispatch(renderDataOverlay(synapse));
            next(loadingSuccess(msg, action.type));
          }, (error) => {
            console.error(msg, error);
            next(raiseError(error));
          });
          break;
        case 'contact':
          contactService.getByID(id).then((contact) => {
            store.dispatch(renderDataOverlay(contact));
            next(loadingSuccess(msg, action.type));
          }, (error) => {
            console.error(msg, error);
            next(raiseError(error));
          });
          break;
        default:
          break;
      }
      break;
    }

    case SET_ORIGINAL_COLORS_FLASHING: {
      const widget = getWidget(store, action.viewerId);
      const instances = setOriginalColorSelectedInstances(
        widget.config.instances,
        action.uids,
      );
      widget.config = {
        ...widget.config,
        instances,
      };
      store.dispatch(updateWidget(widget));
      break;
    }

    case UPDATE_TIMEPOINT_VIEWER: {
      const {
        viewerId, newTimePoint, currTimePoint, widgetType, newViewer,
      } = action;
      const currentWidget = getWidget(store, viewerId);

      const addedObjectsToViewer = currentWidget.config?.addedObjectsToViewer || [];

      if (currTimePoint !== newTimePoint) {
        if (widgetType === VIEWERS.CphateViewer) {
          const msg = 'Updating cphate';
          next(loading(msg, action.type));
          cphateService
            .getCphateByTimepoint(newTimePoint)
            .then((cphate) => {
              if (cphate) {
                const devStages = store.getState().devStages.neuroSCAN;
                const cphateInstances = cphateService.getInstances(cphate, devStages);
                createSimpleInstancesFromInstances(cphateInstances)
                  .then(() => {
                    let newWidget;
                    const foundNewWidget = findWidget(store, newTimePoint, widgetType);
                    if (foundNewWidget) {
                      newWidget = foundNewWidget;
                    }
                    newWidget = createWidget(store, newTimePoint, widgetType);
                    store
                      .dispatch(
                        addToWidget(
                          newWidget,
                          cphateInstances,
                          false,
                        ),
                      );
                    next(loadingSuccess(msg, action.type));
                  });
              }
            }, (error) => {
              console.error(msg, error);
              next(raiseError(msg));
            });
        } else {
          console.debug('Updating viewer timepoint');
          const msg = 'Updating viewer timepoint';
          next(loading(msg, action.type));
          // next(addToWidget(newWidget, [], false, addedObjectsToViewer));
          const neurons = getInstancesOfType(addedObjectsToViewer, NEURON_TYPE) || [];
          const contacts = getInstancesOfType(addedObjectsToViewer, CONTACT_TYPE) || [];
          const synapses = getInstancesOfType(addedObjectsToViewer, SYNAPSE_TYPE) || [];
          const scale = getInstancesOfType(addedObjectsToViewer, SCALE_TYPE) || [];
          // const nerveRing = getInstancesOfType(addedObjectsToViewer, NERVE_RING_TYPE) || [];

          const newNeurons = await fetchDataForEntity(neuronService, newTimePoint, neurons);
          const newContacts = await fetchDataForEntity(contactService, newTimePoint, contacts);
          const newSynapses = await fetchDataForEntity(synapseService, newTimePoint, synapses);
          const newScale = await fetchDataForEntity(scaleService, newTimePoint, scale);
          // const newScale = await scaleService.getByUID(newTimePoint);
          // const newNerveRing = await fetchDataForEntity(nerveRingService,
          // newTimePoint, nerveRing);

          const newInstances = [
            ...newNeurons,
            ...newContacts,
            ...newSynapses,
            ...newScale,
            // ...newNerveRing,
          ].map((i) => mapToInstance(i, store.getState().devStages.neuroSCAN));
          createSimpleInstancesFromInstances(newInstances).then(async () => {
            const newWidget = createWidget(store, newTimePoint, widgetType);
            store.dispatch(addToWidget(newWidget, newInstances, false, addedObjectsToViewer));
            next(loadingSuccess(msg, action.type));
          });
          // cleanEmpty(store);
        }
      }
      break;
    }

    case DELETE_FROM_WIDGET: {
      const widget = getWidget(store, action.viewerId);
      const { addedObjectsToViewer } = widget.config;
      const msg = 'Removing from widget';
      next(loading(msg, action.type));
      store.dispatch(addToWidget(widget, [], true, addedObjectsToViewer));
      next(loadingSuccess(msg, action.type));
      break;
    }

    case ADD_INSTANCES_TO_GROUP: {
      const { viewerId, instances, group } = action;
      const widget = getWidget(store, viewerId);
      // set groupe of instance(s)
      widget.config.instances = updateInstanceGroup(
        widget.config.instances,
        instances,
        group,
      );
      store.dispatch(layoutActions.updateWidget(widget));
      break;
    }

    case SET_INSTANCES_COLOR: {
      const { viewerId, instances, color } = action;
      const widget = getWidget(store, viewerId);
      // set color of instance(s)
      widget.config.instances = setInstancesColor(
        widget.config.instances,
        instances,
        color,
      );
      store.dispatch(layoutActions.updateWidget(widget));
      break;
    }

    case ROTATE_START_ALL: {
      const state = store.getState();
      const newRotateState = cameraControlsRotateState.STARTING;
      Object.values(state.widgets)
        .filter((w) => w.config.rotate === cameraControlsRotateState.STOP)
        .forEach((w) => {
          store.dispatch(updateWidgetConfig(w.config.viewerId, {
            ...w.config,
            rotate: newRotateState,
          }));
        });
      break;
    }

    case ROTATE_STOP_ALL: {
      const state = store.getState();
      const newRotateState = cameraControlsRotateState.STOPPING;
      Object.values(state.widgets)
        .filter((w) => w.config.rotate === cameraControlsRotateState.ROTATING)
        .forEach((w) => {
          store.dispatch(updateWidgetConfig(w.config.viewerId, {
            ...w.config,
            rotate: newRotateState,
          }));
        });
      break;
    }

    case TOGGLE_INSTANCE_HIGHLIGHT: {
      const { viewerId, optionName } = action.payload;
      const state = store.getState();
      const currentHighlighted = state.widgets[viewerId]?.config?.highlightedInstances || [];
      const isCurrentlyHighlighted = currentHighlighted.includes(optionName);
      const updatedHighlighted = isCurrentlyHighlighted
        ? currentHighlighted.filter((name) => name !== optionName)
        : [...currentHighlighted, optionName];
      store.dispatch(updateWidgetConfig(viewerId, {
        ...state.widgets[viewerId].config,
        highlightedInstances: updatedHighlighted,
      }));
      break;
    }

    default:
      next(action);
  }
};

export default middleware;
