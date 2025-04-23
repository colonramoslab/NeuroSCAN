import React from 'react';
import { useSelector } from 'react-redux';
import {
  makeStyles,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Box,
  Button,
  Typography,
  Divider,
} from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';
import CHEVRON from '../../images/chevron-right.svg';
import { resetDataOverlay } from '../../services/instanceHelpers';
import vars from '../../styles/constants';
import { formatSynapseUID } from '../../utilities/functions';

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
        '& .stat-title': {
          fontWeight: 700,
          fontSize: '14px',
          margin: '10px 4px',
        },
      },
    },
  },
}));

const CellStats = ({ dataOverlay }) => (
  <Accordion>
    <AccordionSummary
      expandIcon={<img src={CHEVRON} width="auto" height="auto" alt="CHEVRON" />}
      IconButtonProps={{ disableRipple: true }}
    >
      <Typography variant="h5" class="stat-title">Cell Stats</Typography>
    </AccordionSummary>
    <AccordionDetails>
      <Box className="data-overlay-body">
        {dataOverlay.cell_stats.volume && (
          <p>
            <strong>Volume: </strong>
            {`${Math.round(dataOverlay.cell_stats.volume).toLocaleString()}nm`}
            <sup>3</sup>
          </p>
        )}
        {dataOverlay.cell_stats.surface_area && (
          <p>
            <strong>Cell surface area: </strong>
            {`${Math.round(dataOverlay.cell_stats.surface_area).toLocaleString()}nm`}
            <sup>2</sup>
          </p>
        )}
      </Box>
    </AccordionDetails>
  </Accordion>
);

const PatchStats = ({ dataOverlay }) => (
  <Accordion>
    <AccordionSummary
      expandIcon={<img src={CHEVRON} width="auto" height="auto" alt="CHEVRON" />}
      IconButtonProps={{ disableRipple: true }}
    >
      <Typography variant="h5" class="stat-title">Patch Stats</Typography>
    </AccordionSummary>
    <AccordionDetails>
      <Box className="data-overlay-body">
        {dataOverlay.patch_stats.patch_surface_area && (
          <p>
            <strong>Total surface area: </strong>
            {`${Math.round(dataOverlay.patch_stats.patch_surface_area).toLocaleString()}nm`}
            <sup>2</sup>
            {dataOverlay.cell_stats.surface_area && (
              <>
                &nbsp;(
                {`${((
                  dataOverlay.patch_stats.patch_surface_area
                    / dataOverlay.cell_stats.surface_area
                )
                  * 100
                ).toExponential(2)}%` }
                )
              </>
            )}
          </p>
        )}
      </Box>
    </AccordionDetails>
  </Accordion>
);

const synapseItems = (synapses) => synapses.map((synapse) => (
  <p>
    <strong dangerouslySetInnerHTML={{
      __html: formatSynapseUID(synapse.name),
    }}
    />
    :&nbsp;
    {synapse.count}
  </p>
));

const SynapseStats = ({ dataOverlay }) => (
  <Accordion>
    <AccordionSummary
      expandIcon={<img src={CHEVRON} width="auto" height="auto" alt="CHEVRON" />}
      IconButtonProps={{ disableRipple: true }}
    >
      <Typography variant="h5" class="stat-title">Synapse Stats</Typography>
    </AccordionSummary>
    <AccordionDetails>
      <Box className="data-overlay-body">
        {dataOverlay.volume && (
          <p>
            <strong>Total synapses of type: </strong>
            {`${dataOverlay.synapse_stats.total_type_count}`}
            {dataOverlay.synapse_stats.total_cell_synapse_count && (
              <>
                &nbsp;(
                {`${Math.round(
                  (
                    dataOverlay.synapse_stats.total_type_count
                    / dataOverlay.synapse_stats.total_cell_synapse_count
                  )
                  * 100,
                )}%` }
                )
              </>
            )}
          </p>
        )}
        {synapseItems(dataOverlay.synapse_stats.connections)}
      </Box>
    </AccordionDetails>
  </Accordion>
);

const dataOverlayAccordion = (dataOverlay) => {
  if (!dataOverlay) {
    return <></>;
  }

  if (dataOverlay?.instanceType === 'neuron') {
    return <CellStats dataOverlay={dataOverlay} />;
  }

  if (dataOverlay?.instanceType === 'contact') {
    return (
      <>
        <PatchStats dataOverlay={dataOverlay} />
      </>
    );
  }

  if (dataOverlay?.instanceType === 'synapse') {
    return (
      <>
        <SynapseStats dataOverlay={dataOverlay} />
      </>
    );
  }

  return <></>;
};

const dataOverlayTitle = (dataOverlay) => {
  let title = '';

  if (!dataOverlay) {
    return title;
  }

  const { uid } = dataOverlay;

  if (!uid || uid.length === 0) {
    return title;
  }

  switch (dataOverlay.instanceType) {
    case 'neuron':
      title = dataOverlay.uid;
      break;
    case 'contact':
      title = dataOverlay.uid;
      break;
    case 'synapse':
      title = formatSynapseUID(dataOverlay.uid);
      break;
    default:
      break;
  }

  return title;
};

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
            <Typography component="h3" className="data-overlay-title"><span dangerouslySetInnerHTML={{ __html: dataOverlayTitle(dataOverlay) }} /></Typography>
            <Button onClick={() => resetDataOverlay()} fontSize="large" className="data-overlay-icon">
              <CloseIcon />
            </Button>
          </Box>
          <Divider />
          <Box className="data-overlay-body">
            {dataOverlayAccordion(dataOverlay)}
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
