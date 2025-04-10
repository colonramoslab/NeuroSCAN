// import viewersReducer from './viewers';
import devStagesReducer from './devStages';
import searchReducer from './search';
import miscReducer from './misc';
import promoterReducer from './promoters';
import selectedInstanceToDeleteReducer from './selectedInstanceToDelete';
import dataOverlayReducer from './dataOverlay';

export default {
  // viewers: viewersReducer,
  devStages: devStagesReducer,
  dataOverlay: dataOverlayReducer,
  misc: miscReducer,
  search: searchReducer,
  promoterDB: promoterReducer,
  selectedInstanceToDelete: selectedInstanceToDeleteReducer,
};
