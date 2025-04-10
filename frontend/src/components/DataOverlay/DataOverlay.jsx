import React from 'react';
import { useSelector } from 'react-redux';
import {
  makeStyles,
  Box,
  Drawer,
  Typography,
  IconButton,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@material-ui/core';
import vars from '../../styles/constants';

const useStyles = makeStyles((theme) => ({
  root: {
    position: 'fixed',
    bottom: 0,
    right: 0,
    border: `4px solid ${vars.primaryColor}`,
    backgroundColor: '#ffffff',
  },
}));

const DataOverlay = () => {
  const classes = useStyles();
  const storeData = useSelector((state) => state.dataOverlay);
  const { dataOverlay } = storeData;

  return (
    (dataOverlay?.uid ? (
      <Box className={classes.root}>
        <Box className="data-overlay">
          <Typography>{dataOverlay.uid}</Typography>
        </Box>
      </Box>
    ) : (
      <></>
    ))
  );
};

export default DataOverlay;
