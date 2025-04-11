import React from 'react';
import { useSelector } from 'react-redux';
import {
  makeStyles,
  Box,
  Button,
  Typography,
  Divider,
} from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';
import { resetDataOverlay } from '../../services/instanceHelpers';
import vars from '../../styles/constants';

const useStyles = makeStyles((theme) => ({
  root: {
    position: 'fixed',
    bottom: 0,
    right: 0,
    border: `4px solid ${vars.primaryColor}`,
    backgroundColor: '#ffffff',
    width: '300px',
    maxWidth: '300px',
    maxHeight: '70vh',
    '& .data-overlay': {
      '&-icon': {
        minWidth: 'initial',
      },
      '&-header': {
        display: 'flex',
        flexDirection: 'row',
        justifyContent: 'space-between',
      },
      '&-title': {
        padding: '10px 16px',
        fontWeight: 700,
        overflowX: 'hidden',
        textOverflow: 'ellipsis',
        maxWidth: 'calc(100% - 40px)',
      },
      '&-body': {
        display: 'flex',
        flexDirection: 'row',
        justifyContent: 'space-between',
      },
    },
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
          <Box className="data-overlay-header">
            <Typography component="h3" className="data-overlay-title">{dataOverlay.id}</Typography>
            <Button onClick={() => resetDataOverlay()} fontSize="large" className="data-overlay-icon">
              <CloseIcon />
            </Button>
          </Box>
          <Divider />
          <Box className="data-overlay-body" />
        </Box>
      </Box>
    ) : (
      <></>
    ))
  );
};

export default DataOverlay;
