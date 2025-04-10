export const CLEAR_DATA_OVERLAY = 'CLEAR_DATA_OVERLAY';
export const RENDER_DATA_OVERLAY = 'RENDER_DATA_OVERLAY';

export const clearDataOverlay = (() => ({
  type: CLEAR_DATA_OVERLAY,
}));

export const renderDataOverlay = ((data) => ({
  type: RENDER_DATA_OVERLAY,
  overlayData: data,
}));
