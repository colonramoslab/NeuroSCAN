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
    width: '360px',
    maxWidth: '360px',
    maxHeight: '70vh',
    overflowY: 'scroll',
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
        padding: '8px 0',
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'space-between',
        '& p': {
          margin: '8px 0',
          padding: '0 16px',
        },
      },
    },
  },
}));

const sumSynapses = (synapses) => synapses.reduce((acc, curr) => acc + curr.count, 0);

const DataOverlay = () => {
  const classes = useStyles();
  const storeData = useSelector((state) => state.dataOverlay);
  const { dataOverlay } = storeData;

  return (
    (dataOverlay?.uid ? (
      <Box className={classes.root}>
        <Box className="data-overlay">
          <Box className="data-overlay-header">
            <Typography component="h3" className="data-overlay-title">{dataOverlay.uid}</Typography>
            <Button onClick={() => resetDataOverlay()} fontSize="large" className="data-overlay-icon">
              <CloseIcon />
            </Button>
          </Box>
          <Divider />
          <Box className="data-overlay-body">
            {dataOverlay.volume && (
              <p>
                <strong>Volume: </strong>
                {`${Math.round(dataOverlay.volume).toLocaleString()}nm`}
                <sup>3</sup>
              </p>
            )}
            {dataOverlay.surface_area && (
              <p>
                <strong>Cell surface area: </strong>
                {`${Math.round(dataOverlay.surface_area).toLocaleString()}nm`}
                <sup>2</sup>
              </p>
            )}
            {dataOverlay.contact_surface_area && (
              <p>
                <strong>Contact surface area: </strong>
                {`${Math.round(dataOverlay.contact_surface_area).toLocaleString()}nm`}
                <sup>2</sup>
                &nbsp;
                (
                { `${Math.round(dataOverlay.total_nrc_surface_area).toLocaleString()}nm` }
                <sup>2</sup>
                &nbsp;whole nerve ring
                )
              </p>
            )}
            {dataOverlay.total_contact_surface_area && (
              <p>
                <strong>Patch surface area: </strong>
                {`${Math.round(dataOverlay.total_contact_surface_area).toLocaleString()}nm`}
                <sup>2</sup>
              </p>
            )}
            {dataOverlay.synapses && dataOverlay.synapses.length > 0 && (
              <p>
                <strong>Synapse count: </strong>
                {`${sumSynapses(dataOverlay.synapses)}`}
                &nbsp;(
                { dataOverlay.total_nr_synapses }
                &nbsp;whole nerve ring
                )
              </p>
            )}
          </Box>
        </Box>
      </Box>
    ) : (
      <></>
    ))
  );
};

export default DataOverlay;
