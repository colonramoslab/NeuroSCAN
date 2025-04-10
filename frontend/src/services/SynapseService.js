import qs from 'qs';
import {
  SYNAPSE_TYPE, backendClient, maxRecordsPerFetch, NEURON_TYPE,
} from '../utilities/constants';

const synapsesBackendUrl = '/synapses';

/* eslint class-methods-use-this:
    ["error", { "exceptMethods": ["getByID", "constructQuery", "getByUID"] }]
*/
export class SynapseService {
  async getByID(id) {
    const response = await backendClient.get(`${synapsesBackendUrl}/${id}`);
    return response.data.map((synapse) => ({
      instanceType: SYNAPSE_TYPE,
      ...synapse,
    }));
  }

  async getByUID(timePoint, uids = []) {
    // UIDS can have ampersands in them, so we need to encode them
    const encodedUids = uids.map((uid) => encodeURIComponent(uid));
    const query = `timepoint=${timePoint}${encodedUids.map((uid) => `&uid=${uid}`).join('')}`;
    const response = await backendClient.get(`${synapsesBackendUrl}?${query}`);
    return response.data.map((synapse) => ({
      instanceType: SYNAPSE_TYPE,
      ...synapse,
    }));
  }

  constructQuery(searchState) {
    const { filters } = searchState;
    const { searchTerms, timePoint } = filters;
    const results = searchState.results.synapses;
    let start = 0;

    if (searchState.start !== undefined) {
      start = searchState.start;
    } else if (results.items.length > 0) {
      start = results.items.length;
    }

    const query = {
      timepoint: timePoint,
      start,
      limit: searchState.limit ? searchState.limit : maxRecordsPerFetch,
      sort: 'uid:ASC',
    };

    if (filters.synapsesFilter.chemical) {
      query.type = 'chemical';
    }
    if (filters.synapsesFilter.electrical) {
      query.type = 'electrical';
    }
    if (filters.synapsesFilter.preNeuron) {
      query.pre_neuron = filters.synapsesFilter.preNeuron;
    }
    if (filters.synapsesFilter.postNeuron) {
      query.post_neuron = filters.synapsesFilter.postNeuron;
    }
    let queryString = qs.stringify(query);

    if (searchTerms.length > 0) {
      const searchTermsString = searchTerms.map((term) => `&uid=${term}`).join('');
      queryString += searchTermsString;
    }

    return queryString;
  }

  async search(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${synapsesBackendUrl}?${query}`);
    return response.data.map((synapse) => ({
      instanceType: SYNAPSE_TYPE,
      ...synapse,
    }));
  }

  async getAll(searchState) {
    const query = this.constructQuery({
      ...searchState,
      start: searchState.start,
      limit: searchState.limit,
    });
    const response = await backendClient.get(`${synapsesBackendUrl}?${query}`);
    return response.data.map((synapse) => ({
      instanceType: SYNAPSE_TYPE,
      ...synapse,
    }));
  }

  async totalCount(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${synapsesBackendUrl}/count?${query}`);
    return response.data;
  }
}

export default new SynapseService();
