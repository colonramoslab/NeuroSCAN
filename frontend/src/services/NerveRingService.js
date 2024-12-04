import { NERVE_RING_TYPE, backendClient } from '../utilities/constants';
// eslint-disable-next-line import/no-cycle
import { getLocationPrefixFromType, buildColor } from './instanceHelpers';

const nerveRingsBackendUrl = '/nerve-rings';

/* eslint class-methods-use-this:
    ["error", { "exceptMethods": ["getInstances", "getNerveRingByTimepoint"] }]
*/
export class NerveRingService {
  mapNerveRingInstance(ring, fileUrl, fileType) {
    return {
      id: ring.id,
      uid: ring.uid,
      uidFromDb: ring.uid,
      name: `Nerve Ring ${ring.timepoint}`,
      selected: false,
      instanceType: NERVE_RING_TYPE,
      group: null,
      color: buildColor(ring.color),
      content: {
        type: fileType,
        location: `${fileUrl}${ring.filename}`,
        fileName: ring.filename,
      },
      getId: () => this.id,
    };
  }

  getInstances(nerveRing) {
    const nerveRingFile = getLocationPrefixFromType({
      timepoint: nerveRing.timepoint,
      instanceType: NERVE_RING_TYPE,
    });
    return [this.mapNerveRingInstance(nerveRing, nerveRingFile, 'url')];
  }

  async getNerveRingByTimepoint(timepoint) {
    return backendClient
      .get(nerveRingsBackendUrl, {
        params: {
          timepoint,
        },
      })
      .then((response) => response.data[0]);
  }
}

export default new NerveRingService();
