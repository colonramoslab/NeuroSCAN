import { VIEWERS } from '../../utilities/constants';
// eslint-disable-next-line import/no-cycle
import InstanceViewer from '../viewers/InstanceViewer';

const componentMap = {
  [VIEWERS.InstanceViewer]: InstanceViewer,
  [VIEWERS.CphateViewer]: InstanceViewer,
  [VIEWERS.NerveRingViewer]: InstanceViewer,
};

export default componentMap;
