import React from 'react';
import { useSelector } from 'react-redux';
import {
  Typography,
  Menu,
  MenuItem,
  Divider,
  makeStyles,
} from '@material-ui/core';
import MENU_CHECKMARK_ON from '../../images/menu-checkmark-on.svg';
import PLUS from '../../images/plus-white.svg';
// eslint-disable-next-line import/no-cycle
import { getViewersFromWidgets } from '../../utilities/functions';
import { VIEWERS } from '../../utilities/constants';

const useStyles = makeStyles(() => ({
  mr_8: {
    marginRight: '.5rem',
  },
}));

const AddToViewerMenu = ({
  anchorEl,
  handleClose,
  handleAddToViewer,
  fullMenu = true,
  useAnchorPosition = false,
  anchorPosition = { top: 0, left: 0 },
}) => {
  const classes = useStyles();
  const timePoint = useSelector((state) => state.search.filters.timePoint);
  const widgets = useSelector((state) => state.widgets);
  const viewers = getViewersFromWidgets(widgets);

  function returnMenu() {
    if (viewers.length > 0) {
      return [
        <Typography key="add-to-viewer-text">Add to existing viewer</Typography>,
        viewers.map((viewer) => {
          const instanceName = anchorEl?.parentNode?.textContent.replace('Add to', '');
          let isEnabled = viewer.config.timePoint === timePoint
            && viewer.config.type === VIEWERS.InstanceViewer;
          if (isEnabled) {
            const isInstancePresent = (inst) => {
              let instName = inst.name;
              if (inst.instanceType === 'synapse') {
                instName = inst.name.replace(/<\/?b>/g, '');
              }
              if (instName === instanceName) {
                isEnabled = false;
                return true;
              }
              return false;
            };
            viewer.config.instances.findIndex(isInstancePresent);
          }
          return (
            <MenuItem key={`add-to-viewer-${viewer.id}`} disabled={!isEnabled} onClick={() => handleAddToViewer(viewer.id)}>
              <img src={MENU_CHECKMARK_ON} className={classes.mr_8} alt="MENU_CHECKMARK_ON" />
              {viewer.name}
            </MenuItem>
          );
        }),
        <Divider key="add-to-viewer-divider" />,
      ];
    }
    return [];
  }

  const handleContextMenu = (event) => {
    event.preventDefault();
  };

  return (
    <Menu
      id="addToViewerMenu"
      className="custom-popover dark right"
      anchorEl={anchorEl}
      keepMounted
      onContextMenu={handleContextMenu}
      open={useAnchorPosition ? true : Boolean(anchorEl)}
      anchorReference={useAnchorPosition ? 'anchorPosition' : undefined}
      anchorPosition={useAnchorPosition ? anchorPosition : undefined}
      anchorOrigin={useAnchorPosition ? undefined : { vertical: 'top', horizontal: 'right' }}
      transformOrigin={useAnchorPosition ? undefined : { vertical: 'top', horizontal: 'right' }}
      onClose={handleClose}
      getContentAnchorEl={null}
    >
      { fullMenu ? (
        returnMenu()
      ) : null}
      <MenuItem key="add-to-new-viewer" onClick={() => handleAddToViewer()}>
        <img src={PLUS} className={classes.mr_8} alt="PLUS" />
        Add to New Viewer
      </MenuItem>
    </Menu>
  );
};

export default AddToViewerMenu;
