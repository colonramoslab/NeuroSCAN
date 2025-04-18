/* eslint-disable import/no-cycle */
import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Popover,
} from '@material-ui/core';
import {
  CONTACT_TYPE,
  NEURON_TYPE,
  SYNAPSE_TYPE,
  CPHATE_TYPE,
  VIEWER_MENU,
} from '../../utilities/constants';
import LayersMenu from './ControlMenus/LayersMenu';
import DevStageMenu from './ControlMenus/DevStageMenu';
import ColorPickerMenu from './ControlMenus/ColorPickerMenu';
import {
  getInstancesOfType,
  getInstancesByGroups,
} from '../../services/instanceHelpers';
import { updateTimePointViewer, deleteFromWidget } from '../../redux/actions/widget';
import WarningModal from '../WarningModal';

const MenuControl = ({
  anchorEl, setAnchorEl, handleClose, open, id, selection, viewerId,
}) => {
  const dispatch = useDispatch();
  const widgets = useSelector((state) => state.widgets);
  const currentWidget = widgets[viewerId];

  const [content, setContent] = useState(null);
  const [timePoint, setTimePoint] = useState(currentWidget?.config?.timePoint || 0);

  const [openWarningModal, setOpenWarningModal] = useState(false);
  const [lostInstances, setLostInstances] = useState([]);

  // const layersList = ['Worm Body', 'Pharynx', 'NerveRing'];
  const downloadFiles = (option) => {
    // console.log(`selected option: ${option}`);
    handleClose();
  };
  let instances = [];
  let addedObjectsToViewer = [];
  if (currentWidget) {
    instances = currentWidget.config.instances;
    addedObjectsToViewer = currentWidget.config.addedObjectsToViewer;
  }
  const groups = getInstancesByGroups(instances);
  const neurons = getInstancesOfType(instances, NEURON_TYPE) || [];
  const contacts = getInstancesOfType(instances, CONTACT_TYPE) || [];
  const synapses = getInstancesOfType(instances, SYNAPSE_TYPE) || [];
  const clusters = getInstancesOfType(instances, CPHATE_TYPE) || [];

  useEffect(() => {
    if (currentWidget && timePoint !== currentWidget?.config?.timePoint) {
      /* TODO: this below is just an hack, it requires a new geppetto
       * version but we are not supporting this anymore. */
      // dispatch(deleteFromWidget(viewerId));
      // This below is the only line should stay in this condition
      dispatch(updateTimePointViewer(
        viewerId,
        timePoint,
        currentWidget.config.timePoint,
        currentWidget.component,
        true,
      ));
      setAnchorEl(null);
      // end of the hack
    }
  }, [timePoint]);

  useEffect(() => {
    if (currentWidget && timePoint !== currentWidget?.timePoint && instances.length !== 0) {
      const lostInstancesArray = instances?.filter((item1) => !addedObjectsToViewer
        .some((item2) => item1.name.toLowerCase() === item2.name.toLowerCase()));
      setLostInstances(lostInstancesArray);
    }
  }, [timePoint, addedObjectsToViewer, instances]);

  useEffect(() => {
    // if (currentWidget
    // && timePoint !== currentWidget?.timePoint
    // && lostInstances.length !== 0) {
    // console.log(lostInstances);
    // const delay = setTimeout(() => {
    // setOpenWarningModal(true);
    // clearTimeout(delay);
    // }, 1000);
    // } else {
    setOpenWarningModal(false);
    // }
  }, [lostInstances]);

  useEffect(() => {
    switch (selection) {
      case VIEWER_MENU.devStage: setContent(
        <DevStageMenu
          timePoint={currentWidget?.config?.timePoint}
          setTimePoint={setTimePoint}
        />,
      );
        break;
      // case VIEWER_MENU.layers: setContent(<LayersMenu layers={layersList} />);
      //   break;
      case VIEWER_MENU.colorPicker:
        setContent(
          <ColorPickerMenu
            dispatch={dispatch}
            viewerId={viewerId}
            groups={groups}
            neurons={neurons}
            contacts={contacts}
            synapses={synapses}
            clusters={clusters}
          />,
        );
        break;
      default:
        setContent(null);
    }
  }, [selection, instances]);

  return (
    <>
      <Popover
        id={id}
        className="custom-popover"
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'left',
        }}
      >
        { content }
      </Popover>
      {
          lostInstances?.length !== 0 && openWarningModal && (
          <WarningModal
            open={openWarningModal}
            handleClose={() => setOpenWarningModal(false)}
            instances={lostInstances}
          />
          )
      }

    </>

  );
};

export default MenuControl;
