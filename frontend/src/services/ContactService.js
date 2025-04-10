import qs from 'qs';
import {
  CONTACT_TYPE, backendClient, maxRecordsPerFetch, SYNAPSE_TYPE,
} from '../utilities/constants';

const contactsUrl = '/contacts';

/* eslint class-methods-use-this:
    ["error", { "exceptMethods": ["getByID", "constructQuery", "getByUID"] }]
*/
export class ContactService {
  async getByID(id) {
    const response = await backendClient.get(`${contactsUrl}/${id}`);
    return {
      ...response.data,
      instanceType: CONTACT_TYPE,
    };
  }

  async getByUID(timePoint, uids = []) {
    const query = `timepoint=${timePoint}${uids.map((uid, i) => `${(i === 0) ? '&' : ''}uid_in=${uid}`).join('&')}`;
    const response = await backendClient.get(`${contactsUrl}?${query}`);
    return response.data.map((contact) => ({
      instanceType: CONTACT_TYPE,
      ...contact,
    }));
  }

  async search(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${contactsUrl}?${query}`);
    return response.data.map((contact) => ({
      instanceType: CONTACT_TYPE,
      ...contact,
    }));
  }

  async getAll(searchState) {
    const query = this.constructQuery({
      ...searchState,
      start: searchState.start,
      limit: searchState.limit,
    });
    const response = await backendClient.get(`${contactsUrl}?${query}`);
    return response.data.map((contact) => ({
      instanceType: CONTACT_TYPE,
      ...contact,
    }));
  }

  async totalCount(searchState) {
    const query = this.constructQuery(searchState);
    const response = await backendClient.get(`${contactsUrl}/count?${query}`);
    return response.data;
  }

  constructQuery(searchState) {
    const { searchTerms, timePoint } = searchState.filters;
    const results = searchState.results.contacts;
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
}

export default new ContactService();
