import { RENDER_DATA_OVERLAY, CLEAR_DATA_OVERLAY } from '../actions/dataOverlay';

export const DEFAULT_DATA_OVERLAY = {};

export default (state = DEFAULT_DATA_OVERLAY, action) => {
  switch (action.type) {
    case RENDER_DATA_OVERLAY:
      return {
        ...state,
        dataOverlay: action.overlayData,
      };
    case CLEAR_DATA_OVERLAY:
      return {
        ...state,
        dataOverlay: DEFAULT_DATA_OVERLAY,
      };
    default:
      return state;
  }
};
