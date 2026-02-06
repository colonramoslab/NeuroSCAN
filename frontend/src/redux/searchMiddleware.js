import * as search from './actions/search';
import doSearch from '../services/helpers';
import doGetAll from '../services/getAllHelper';
import nerveRingService from '../services/NerveRingService';
import cphateService from '../services/CphateService';
import scaleService from '../services/ScaleService';
import {
  ADD_CPHATE, ADD_NERVE_RING, ADD_SCALE, addInstances,
} from './actions/widget';
import { raiseError, loading, loadingSuccess } from './actions/misc';
import { VIEWERS } from '../utilities/constants';

const searchMiddleware = (store) => (next) => (action) => {
  switch (action.type) {
    case search.GET_ALL: {
      const { entity } = action.data;
      next(action);
      const state = store.getState();
      doGetAll(store.dispatch, { ...state.search }, [entity]);
      break;
    }
    case search.UPDATE_FILTERS: {
      next(action);
      const state = store.getState();
      state.search.filters.timePoint = action.timePoint;
      doSearch(store.dispatch, state.search);
      break;
    }

    case search.LOAD_MORE: {
      const { entity } = action.data;
      next({
        type: action.type,
      });
      const state = store.getState();
      doSearch(store.dispatch, state.search, [entity]);
      break;
    }

    case ADD_CPHATE: {
      const { timePoint } = action;
      const msg = 'Add cphate';
      next(loading(msg, action.type));
      cphateService
        .getCphateByTimepoint(timePoint)
        .then((cphate) => {
          if (cphate) {
            const devStages = store.getState().devStages.neuroSCAN;
            const cphateInstances = cphateService.getInstances(cphate, devStages);
            store.dispatch(addInstances(null, cphateInstances, VIEWERS.CphateViewer));
          }
          next(loadingSuccess(msg, action.type));
        }, (error) => {
          console.error(msg, error);
          next(raiseError(msg));
        });
      break;
    }

    case ADD_NERVE_RING: {
      const { timePoint, viewer } = action;
      const viewerId = viewer.id || null;
      if (viewer?.config?.type?.toLowerCase() === 'cphate') {
        return;
      }
      const msg = 'Add nerve ring';
      next(loading(msg, action.type));
      nerveRingService
        .getNerveRingByTimepoint(timePoint)
        .then((ring) => {
          if (ring) {
            const devStages = store.getState().devStages.neuroSCAN;
            const ringInstances = nerveRingService.getInstances(ring, devStages);
            store.dispatch(addInstances(viewerId, ringInstances, VIEWERS.InstanceViewer));
          }
          next(loadingSuccess(msg, action.type));
        },
        (e) => {
          console.error(e);
          next(raiseError(msg));
        });
      break;
    }

    case ADD_SCALE: {
      const { timePoint } = action;
      const msg = 'Add scale';
      next(loading(msg, action.type));
      scaleService
        .getScaleByTimepoint(timePoint)
        .then((scale) => {
          if (scale) {
            const devStages = store.getState().devStages.neuroSCAN;
            const scaleInstances = scaleService.getInstances(scale, devStages);
            store.dispatch(addInstances(null, scaleInstances, VIEWERS.InstanceViewer));
          }
          next(loadingSuccess(msg, action.type));
        },
        (error) => {
          console.error(msg, error);
          next(raiseError(msg));
        });
      break;
    }

    default:
      next(action);
  }
};

export default searchMiddleware;
