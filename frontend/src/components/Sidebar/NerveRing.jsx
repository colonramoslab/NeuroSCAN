import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Typography, Box } from '@material-ui/core';
import AddToViewerMenu from './AddToViewerMenu';
import { addNerveRing } from '../../redux/actions/widget';
import { getViewersFromWidgets } from '../../utilities/functions';
import PLUS from '../../images/plus.svg';

const NerveRing = ({ timePoint }) => {
  const dispatch = useDispatch();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [menuOpen, setMenuOpen] = React.useState(false);
  const widgets = useSelector((state) => state.widgets);

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
    setMenuOpen(true);
  };

  const handleClose = () => {
    setAnchorEl(null);
    setMenuOpen(false);
  };

  const handleAddToViewer = async (viewerId = null) => {
    const viewers = getViewersFromWidgets(widgets);
    let viewer = null;

    if (viewers.length > 0) {
      viewer = viewers.find((v) => v.id === viewerId);
    }

    handleClose();
    dispatch(addNerveRing(viewer.config.timePoint, viewer));
  };

  return (
    <Box className="wrap" id="nerveRing-id">
      <Typography component="h5" onClick={handleClick}>
        <img src={PLUS} alt="Plus" />
        Add Nerve Ring (add last)
      </Typography>
      {menuOpen && (
        <AddToViewerMenu
          handleClose={handleClose}
          handleAddToViewer={handleAddToViewer}
          anchorEl={anchorEl}
          kind="nervering"
        />
      )}
    </Box>
  );
};

export default NerveRing;
