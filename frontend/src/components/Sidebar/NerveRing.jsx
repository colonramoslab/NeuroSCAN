import React from 'react';
import { useDispatch } from 'react-redux';
import { Typography, Box } from '@material-ui/core';
import { addNerveRing } from '../../redux/actions/widget';
import PLUS from '../../images/plus.svg';

const NerveRing = ({ timePoint }) => {
  const dispatch = useDispatch();

  const createNerveRingViewer = async () => {
    dispatch(addNerveRing(timePoint));
  };

  return (
    <Box className="wrap" onClick={createNerveRingViewer} id="nerveRing-id">
      <Typography component="h5">
        <img src={PLUS} alt="Plus" />
        Add Nerve Ring
      </Typography>
    </Box>
  );
};

export default NerveRing;
