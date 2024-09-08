import qs from 'qs';
import { NERVE_RING_TYPE, backendClient, maxRecordsPerFetch } from '../utilities/constants';

const nerveRingsBackendUrl = '/nerve-rings';

/* eslint class-methods-use-this:
    ["error", { "exceptMethods": ["getById", "search", "constructQuery", "getByUID"] }]
*/
export class NerveRingService {
  async getByUID(timePoint, uids = []) {
    const query = `timepoint=${timePoint}${uids.map((uid) => `&uid_in=${uid}`).join('')}`;
    const response = await backendClient.get(`${nerveRingsBackendUrl}?${query}`);
    return response.data.map((nerveRing) => ({
      instanceType: NERVE_RING_TYPE,
      ...nerveRing,
    }));
  }

  constructQuery(searchState) {
    const { searchTerms, timePoint } = searchState.filters;
    const results = searchState.results.nerverings;
    return qs.stringify({
      _where: [
        { timepoint: timePoint },
        // { _or: searchTerms.map((term) => ({ uid_contains: term })) },
      ],
      _sort: 'uid:ASC',
      _start: searchState?.limit ? searchState.start : results.items.length,
      _limit: searchState?.limit || maxRecordsPerFetch,
    });
  }

  async search(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${nerveRingsBackendUrl}?${query}`);
    return response.data.map((nerveRing) => ({
      instanceType: NERVE_RING_TYPE,
      ...nerveRing,
    }));
  }

  async getAll(searchState) {
    const query = this.constructQuery({
      ...searchState,
      start: searchState.start,
      limit: searchState.limit,
    });
    const response = await backendClient.get(`${nerveRingsBackendUrl}?${query}`);
    return response.data.map((nerveRing) => ({
      instanceType: NERVE_RING_TYPE,
      ...nerveRing,
    }));
  }

  async totalCount(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${nerveRingsBackendUrl}/count?${query}`);
    return response.data;
  }
}

export default new NerveRingService();
