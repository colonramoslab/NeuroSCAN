import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import './App.css';
import {
  BrowserRouter as Router,
  Switch,
  Route,
} from 'react-router-dom';
import { MuiThemeProvider } from '@material-ui/core';
import CssBaseline from '@material-ui/core/CssBaseline';
import Loader from '@metacell/geppetto-meta-ui/loader/Loader';
import posthog from 'posthog-js';
import NeuroScan from './pages/NeuroScan';
import PromoterDB from './pages/PromoterDB';
import About from './pages/About';
import theme from './theme';
import { addDevStages } from './redux/actions/devStages';
import { loadPromoters } from './redux/actions/promoters';
import '@metacell/geppetto-meta-ui/flex-layout/style/dark.scss';
import './flexLayout.css';
import { CANVAS_STARTED } from './utilities/constants';

if (window.location.host === 'promoters.wormguides.org' || window.location.host === 'neuroscan.net') {
  posthog.init('phc_cuoRrQDAJFfkKureV66PysaW5hkbrmsuzC4A1SRDAtk', {
    api_host: 'https://us.i.posthog.com',
  });
}

const Manager = require('@metacell/geppetto-meta-client/common/Manager').default;

const GEPPETTO = {};
window.GEPPETTO = GEPPETTO;
GEPPETTO.Resources = require('@metacell/geppetto-meta-core/Resources').default;

GEPPETTO.Manager = new Manager();
window.Instances = [];

const App = () => {
  const dispatch = useDispatch();
  const misc = useSelector((state) => state.misc);
  const loadingMessages = ['Loading...'];

  useEffect(() => {
    dispatch(addDevStages());
    dispatch(loadPromoters());
  }, []);

  const getActiveStatus = () => {
    const loading = Object.keys(misc?.loading || {}).length > 0;
    const canvasLoading = misc?.canvas === CANVAS_STARTED;
    if (loading || canvasLoading) {
      return true;
    }
    return false;
  };

  let mainComponent = NeuroScan;
  const isPromoterDB = window.location.host === 'promoterdb.bilte.co' || window.location.host === 'promoters.wormguides.org' || window.location.search.includes('promoterdb');
  if (isPromoterDB) {
    mainComponent = PromoterDB;
  }

  return (
    <MuiThemeProvider theme={theme}>
      <CssBaseline />
      <Loader active={getActiveStatus()} messages={loadingMessages} className="custom-loading-text" />
      <Router>
        <Switch>
          <Route exact path="/" component={mainComponent} />
          <Route path="/about">
            <About />
          </Route>
          {/* <Route exact path="/promoter" component={PromoterDB} /> */}
        </Switch>
      </Router>
    </MuiThemeProvider>
  );
};

export default App;
