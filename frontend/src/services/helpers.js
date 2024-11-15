import neuronService from './NeuronService';
import contactService from './ContactService';
import synapseService from './SynapseService';
import nerveRingService from './NerveRingService';
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

const doSearchNerveRings = async (dispatch, searchState) => {
  nerveRingService.totalCount(searchState).then((count) => {
    dispatch(
      search.updateCounters({
        nerverings: count,
      }),
    );
  });
  nerveRingService.search(searchState).then((data) => {
    dispatch(
      search.updateResults({
        nerverings: {
          ...searchState.results.nerverings,
          items: searchState.results.nerverings.items.concat(data),
        },
      }),
    );
  });
};

export default async (dispatch, searchState, entities = ['neurons', 'contacts', 'synapses', 'nerverings']) => {
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
      case 'nerverings': {
        doSearchNerveRings(dispatch, searchState);
        break;
      }

      default:
        break;
    }
  });
};
