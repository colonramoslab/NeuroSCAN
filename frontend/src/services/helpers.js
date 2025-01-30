import neuronService from './NeuronService';
import contactService from './ContactService';
import synapseService from './SynapseService';
// eslint-disable-next-line import/no-cycle
import scaleService from './ScaleService';
// eslint-disable-next-line import/no-cycle
import cphateService from './CphateService';
import * as search from '../redux/actions/search';
//

const doSearchNeurons = async (dispatch, searchState) => {
  neuronService.totalCount(searchState).then((count) => {
    dispatch(
      search.updateCounters({
        neurons: count,
      }),
    );
  });
  neuronService.search(searchState).then((data) => {
    dispatch(
      search.updateResults({
        neurons: {
          ...searchState.results.neurons,
          items: searchState.results.neurons.items.concat(data),
        },
      }),
    );
  });
};

const doSearchSynapses = async (dispatch, searchState) => {
  synapseService.totalCount(searchState).then((count) => {
    dispatch(
      search.updateCounters({
        synapses: count,
      }),
    );
  });
  synapseService.search(searchState).then((data) => {
    dispatch(
      search.updateResults({
        synapses: {
          ...searchState.results.synapses,
          items: searchState.results.synapses.items.concat(data),
        },
      }),
    );
  });
};

const doSearchContacts = async (dispatch, searchState) => {
  contactService.totalCount(searchState).then((count) => {
    dispatch(
      search.updateCounters({
        contacts: count,
      }),
    );
  });
  contactService.search(searchState).then((data) => {
    dispatch(
      search.updateResults({
        contacts: {
          ...searchState.results.contacts,
          items: searchState.results.contacts.items.concat(data),
        },
      }),
    );
  });
};

const doSearchCphate = async (dispatch, searchState) => {
  const { timePoint } = searchState.filters;
  cphateService.totalCount(timePoint).then((count) => {
    dispatch(
      search.updateCounters({
        cphate: count,
      }),
    );
  });
  cphateService.getCphateByTimepoint(timePoint).then((data) => {
    dispatch(
      search.updateResults({
        cphate: {
          ...searchState.results.cphate,
          items: searchState.results.cphate.items.concat(data),
        },
      }),
    );
  });
};

const doSearchScale = async (dispatch, searchState) => {
  const { timePoint } = searchState.filters;
  dispatch(
    search.updateCounters({
      scale: 1,
    }),
  );
  scaleService.getScaleByTimepoint(timePoint).then((data) => {
    dispatch(
      search.updateResults({
        scale: {
          items: [data],
        },
      }),
    );
  });
};

export default async (dispatch, searchState, entities = ['neurons', 'contacts', 'synapses', 'scale', 'cphate']) => {
  entities.forEach((entity) => {
    switch (entity) {
      case 'neurons': {
        doSearchNeurons(dispatch, searchState);
        break;
      }
      case 'contacts': {
        doSearchContacts(dispatch, searchState);
        break;
      }
      case 'synapses': {
        doSearchSynapses(dispatch, searchState);
        break;
      }
      case 'scale': {
        doSearchScale(dispatch, searchState);
        break;
      }
      case 'cphate': {
        doSearchCphate(dispatch, searchState);
        break;
      }

      default:
        break;
    }
  });
};
