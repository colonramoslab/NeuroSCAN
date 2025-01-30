import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Typography, Box } from '@material-ui/core';
import { addNerveRing } from '../../redux/actions/widget';
import { getViewersFromWidgets } from '../../utilities/functions';
import PLUS from '../../images/plus.svg';

const NerveRing = ({ timePoint }) => {
  const dispatch = useDispatch();
  const widgets = useSelector((state) => state.widgets);

  const createNerveRingViewer = async () => {
    const viewers = getViewersFromWidgets(widgets);
    let currentViewer = null;

    if (viewers.length > 0) {
      currentViewer = viewers.find((viewer) => viewer.config.timePoint === timePoint);
    }

    dispatch(addNerveRing(timePoint, currentViewer));
  };

  return (
    <Box className="wrap" onClick={createNerveRingViewer} id="nerveRing-id">
      <Typography component="h5">
        <img src={PLUS} alt="Plus" />
        Add Nerve Ring (add last)
      </Typography>
    </Box>
  );
};

export default NerveRing;
