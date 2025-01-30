import qs from 'qs';
import {
  SYNAPSE_TYPE, backendClient, maxRecordsPerFetch, NEURON_TYPE,
} from '../utilities/constants';

const synapsesBackendUrl = '/synapses';

/* eslint class-methods-use-this:
    ["error", { "exceptMethods": ["getById", "constructQuery", "getByUID"] }]
*/
export class SynapseService {
  async getById(id) {
    return {
      id,
      uid: 'TestUIDContact1',
      content: {
        type: 'url',
        location: 'https://raw.githubusercontent.com/MetaCell/geppetto-meta/development/geppetto.js/geppetto-ui/src/3d-canvas/showcase/examples/Sketch_Volume_Viewer_AIB_Rby_AIAR_AIB_Rby_AIAR_1_1_0000_green_0_24947b6670.gltf',
        fileName: 'Sketch_Volume_Viewer_AIB_Rby_AIAR_AIB_Rby_AIAR_1_1_0000_green_0_24947b6670.gltf',
      },
      getId: () => this.id,
    };
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
