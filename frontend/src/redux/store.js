import { createStore } from '@metacell/geppetto-meta-client/common';
import baseLayout from '../components/layout/defaultLayout';
// eslint-disable-next-line import/no-cycle
import componentMap from '../components/layout/componentMap';
import reducers from './reducers/all';
// eslint-disable-next-line import/no-cycle
import middleware from './middleware';
// eslint-disable-next-line import/no-cycle
import searchMiddleware from './searchMiddleware';
import promotersSearchMiddleware from './promotersSearchMiddleware';

import { SEARCH_DEFAULT_STATUS } from './reducers/search';
import { DEVSTAGES_DEFAULT_STATUS } from './reducers/devStages';
import { PROMOTERS_DEFAULT_STATUS } from './reducers/promoters';
import { MISC_DEFAULT_STATUS } from './reducers/misc';
import { DEFAULT_DATA_OVERLAY } from './reducers/dataOverlay';

const INIT_STATE = {
  devStages: DEVSTAGES_DEFAULT_STATUS,
  misc: MISC_DEFAULT_STATUS,
  search: SEARCH_DEFAULT_STATUS.search,
  promoterDB: PROMOTERS_DEFAULT_STATUS,
  dataOverlay: DEFAULT_DATA_OVERLAY,
  selectedInstanceToDelete: {},
};

const store = createStore(
  reducers,
  INIT_STATE,
  [middleware, searchMiddleware, promotersSearchMiddleware],
  { baseLayout, componentMap },
);

export default store;
