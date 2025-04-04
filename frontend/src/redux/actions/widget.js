export const ADD_INSTANCES = 'ADD_INSTANCES';
export const CLONE_VIEWER_WITH_INSTANCES_LIST = 'CLONE_VIEWER_CLONE_INSTANCES_LIST';
export const ADD_CPHATE = 'ADD_CPHATE';
export const ADD_NERVE_RING = 'ADD_NERVE_RING';
export const ADD_SCALE = 'ADD_SCALE';
export const ADD_INSTANCES_TO_GROUP = 'ADD_INSTANCES_TO_GROUP';
export const SET_INSTANCES_COLOR = 'SET_INSTANCES_COLOR';
export const UPDATE_TIMEPOINT_VIEWER = 'UPDATE_TIMEPOINT_VIEWER';
export const UPDATE_BACKGROUND_COLOR_VIEWER = 'UPDATE_BACKGROUND_COLOR_VIEWER';
export const UPDATE_WIDGET_CONFIG = 'UPDATE_WIDGET_CONFIG';
export const UPDATE_VIEWER_ROTATE = 'UPDATE_VIEWER_ROTATE';
export const ROTATE_START_ALL = 'ROTATE_START_ALL';
export const ROTATE_STOP_ALL = 'ROTATE_STOP_ALL';
export const INVERT_COLORS_FLASHING = 'INVERT_COLORS_FLASHING';
export const DARKEN_COLORS_FLASHING = 'DARKEN_COLORS_FLASHING';
export const SET_ORIGINAL_COLORS_FLASHING = 'SET_ORIGINAL_COLORS_FLASHING';
export const TOGGLE_INSTANCE_HIGHLIGHT = 'TOGGLE_INSTANCE_HIGHLIGHT';
export const ADD_LAST_SELECTED_INSTANCE = 'ADD_LAST_SELECTED_INSTANCE';
export const DELETE_FROM_WIDGET = 'DELETE_FROM_WIDGET';

export const addInstances = ((viewerId, instances, viewerType = null) => ({
  type: ADD_INSTANCES,
  viewerId,
  viewerType,
  instances,
}));

export const cloneViewerWithInstancesList = ((fromViewerId, instances) => ({
  type: CLONE_VIEWER_WITH_INSTANCES_LIST,
  fromViewerId,
  instances,
}));

export const addCphate = ((timePoint) => ({
  type: ADD_CPHATE,
  timePoint,
}));

export const addNerveRing = (timePoint, viewer) => ({
  type: ADD_NERVE_RING,
  timePoint,
  viewer,
});

export const addScale = (timePoint, viewer) => ({
  type: ADD_SCALE,
  timePoint,
  viewer,
});

export const addInstancesToGroup = ((viewerId, instances, group) => ({
  type: ADD_INSTANCES_TO_GROUP,
  viewerId,
  instances,
  group,
}));

export const setInstancesColor = ((viewerId, instances, color) => ({
  type: SET_INSTANCES_COLOR,
  viewerId,
  instances,
  color,
}));

export const updateTimePointViewer = ((
  viewerId,
  newTimePoint,
  currTimePoint,
  widgetType,
  newViewer,
) => ({
  type: UPDATE_TIMEPOINT_VIEWER,
  viewerId,
  newTimePoint,
  currTimePoint,
  widgetType,
  newViewer,
}));

export const updateBackgroundColorViewer = ((viewerId, backgroundColor) => ({
  type: UPDATE_BACKGROUND_COLOR_VIEWER,
  viewerId,
  backgroundColor,
}));

export const updateWidgetConfig = ((viewerId, config) => ({
  type: UPDATE_WIDGET_CONFIG,
  viewerId,
  config,
}));

export const rotateStartAll = (() => ({
  type: ROTATE_START_ALL,
}));

export const rotateStopAll = (() => ({
  type: ROTATE_STOP_ALL,
}));

export const invertColorsFlashing = ((viewerId, uids) => ({
  type: INVERT_COLORS_FLASHING,
  viewerId,
  uids,
}));

export const darkenColorsFlashing = ((viewerId, uids) => ({
  type: DARKEN_COLORS_FLASHING,
  viewerId,
  uids,
}));

export const setOriginalColors = ((viewerId, uids) => ({
  type: SET_ORIGINAL_COLORS_FLASHING,
  viewerId,
  uids,
}));

export const toggleInstanceHighlight = (viewerId, optionName) => ({
  type: TOGGLE_INSTANCE_HIGHLIGHT,
  payload: { viewerId, optionName },
});

export const addLastSelectedInstance = ((viewerId, uid) => ({
  type: ADD_LAST_SELECTED_INSTANCE,
  viewerId,
  uid,
}));

export const deleteFromWidget = ((viewerId, uids) => ({
  type: DELETE_FROM_WIDGET,
  viewerId,
}));
