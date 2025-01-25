import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Typography, Box } from '@material-ui/core';
import { addCphate } from '../../redux/actions/widget';
import PLUS from '../../images/plus.svg';

const CPhatePlot = ({ timePoint }) => {
  const dispatch = useDispatch();
  const count = useSelector((state) => state.search.counters.cphate);
  const [elClass, setElClass] = useState('');
  const [clickEnabled, setClickEnabled] = useState(true);

  useEffect(() => {
    if (count > 0) {
      setElClass('');
      setClickEnabled(true);
    } else {
      setElClass('disabled');
      setClickEnabled(false);
    }
  }, [count]);

  const createCphateViewer = async () => {
    dispatch(addCphate(timePoint));
  };

  return (
    <Box className="wrap" onClick={clickEnabled ? createCphateViewer : null} id="cphate-id">
      <Typography component="h5" className={elClass}>
        <img src={PLUS} alt="Plus" />
        Add CPHATE Plot
      </Typography>
    </Box>
  );
};

export default CPhatePlot;
