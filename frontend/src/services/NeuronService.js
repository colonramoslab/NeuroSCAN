import qs from 'qs';
import { NEURON_TYPE, backendClient, maxRecordsPerFetch } from '../utilities/constants';

const neuronsBackendUrl = '/neurons';

/* eslint class-methods-use-this:
    ["error", { "exceptMethods": ["getById", "search", "constructQuery", "getByUID"] }]
*/
export class NeuronService {
  async getById(id) {
    try {
      const response = await backendClient.get(`${neuronsBackendUrl}/${id}`);
      return response.data;
    } catch (error) {
      return {
        id,
        uid: 'SAAVR',
        content: {
          type: 'url',
          location: 'https://raw.githubusercontent.com/MetaCell/geppetto-meta/master/geppetto.js/geppetto-ui/src/3d-canvas/showcase/examples/SketchVolumeViewer_SAAVR_SAAVR_1_1_0000_draco.gltf',
          fileName: 'SketchVolumeViewer_SAAVR_SAAVR_1_1_0000_draco.gltf',
        },
        getId: () => this.id,
      };
    }
  }

  async getByUID(timePoint, uids = []) {
    const query = `timepoint=${timePoint}&limit=${uids.length + 1}${uids.map((uid) => `&uid=${uid}`).join('')}`;
    const response = await backendClient.get(`${neuronsBackendUrl}?${query}`);
    return response.data.map((neuron) => ({
      instanceType: NEURON_TYPE,
      ...neuron,
    }));
  }

  constructQuery(searchState) {
    const { searchTerms, timePoint } = searchState.filters;
    const results = searchState.results.neurons;
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

    let queryString = qs.stringify(query);

    if (searchTerms.length > 0) {
      const searchTermsString = searchTerms.map((term) => `&uid=${term}`).join('');
      queryString += searchTermsString;
    }

    return queryString;
  }

  async search(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${neuronsBackendUrl}?${query}`);
    return response.data.map((neuron) => ({
      instanceType: NEURON_TYPE,
      ...neuron,
    }));
  }

  async getAll(searchState) {
    const query = this.constructQuery({
      ...searchState,
      start: searchState.start,
    });
    const response = await backendClient.get(`${neuronsBackendUrl}?${query}`);
    return response.data.map((neuron) => ({
      instanceType: NEURON_TYPE,
      ...neuron,
    }));
  }

  async totalCount(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${neuronsBackendUrl}/count?${query}`);
    return response.data;
  }
}

export default new NeuronService();
