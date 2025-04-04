import { SCALE_TYPE, backendClient } from '../utilities/constants';
// eslint-disable-next-line import/no-cycle
import { getLocationPrefixFromType, buildColor } from './instanceHelpers';

const scaleBackendUrl = '/scales';

/* eslint class-methods-use-this:
    ["error", { "exceptMethods": ["getInstances", "getScaleByTimepoint", "getByUID"] }]
*/
export class ScaleService {
  mapScaleInstance(scale, fileUrl, fileType) {
    return {
      id: scale.id,
      uid: scale.uid,
      uidFromDb: scale.uid,
      name: `${scale.uid}`,
      selected: false,
      instanceType: SCALE_TYPE,
      group: null,
      color: buildColor(scale.color),
      content: {
        type: fileType,
        location: `${fileUrl}${scale.filename}`,
        fileName: scale.filename,
      },
      getId: () => this.id,
    };
  }

  getInstances(scale) {
    const scaleFile = getLocationPrefixFromType({
      timepoint: scale.timepoint,
      instanceType: SCALE_TYPE,
    });
    return this.mapScaleInstance(scale, scaleFile, 'url');
  }

  async getScaleByTimepoint(timepoint) {
    return backendClient
      .get(scaleBackendUrl, {
        params: {
          timepoint,
        },
      })
      .then((response) => response.data.map((scale) => ({
        instanceType: SCALE_TYPE,
        ...scale,
      })));
  }

  async getByUID(timePoint) {
    const query = `timepoint=${timePoint}`;
    const response = await backendClient.get(`${scaleBackendUrl}?${query}`);
    return response.data.map((scale) => ({
      instanceType: SCALE_TYPE,
      ...scale,
    }));
  }
}

export default new ScaleService();
