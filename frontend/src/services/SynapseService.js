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
    const query = `timepoint=${timePoint}${encodedUids.map((uid) => `&uid_in=${uid}`).join('')}`;
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
    const andPart = [];
    andPart.push({ timepoint: timePoint });
    if (searchTerms.length > 0) {
      andPart.push({ searchTerms: searchTerms.join('|') });
    }
    if (filters.synapsesFilter.chemical) {
      andPart.push({
        type: 'chemical',
      });
    }
    if (filters.synapsesFilter.electrical) {
      andPart.push({
        type: 'electrical',
      });
    }
    if (filters.synapsesFilter.preNeuron) {
      andPart.push({
        neuronPre: filters.synapsesFilter.preNeuron,
      });
    }
    if (filters.synapsesFilter.postNeuron) {
      andPart.push({
        postNeuron: filters.synapsesFilter.postNeuron,
      });
    }
    return qs.stringify({
      _where: andPart,
      _sort: 'uid:ASC',
      _start: searchState?.limit ? searchState?.start : results.items.length,
      _limit: searchState?.limit || maxRecordsPerFetch,
    });
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
