import React from 'react';
import { useSelector } from 'react-redux';
import {
  makeStyles,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  Box,
  Button,
  Tooltip,
  Typography,
  Divider,
} from '@material-ui/core';
import CloseIcon from '@material-ui/icons/Close';
import HelpOutlineIcon from '@material-ui/icons/HelpOutline';
import CHEVRON from '../../images/chevron-right.svg';
import HTMLTooltip from '../HTMLTooltip';
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
        // fontWeight: 700,
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
  tooltip: {
    marginLeft: '4px',
  },
}));

const CellStats = ({ dataOverlay }) => (
  <Accordion>
    <AccordionSummary
      expandIcon={
        <img src={CHEVRON} width="auto" height="auto" alt="CHEVRON" />
      }
      IconButtonProps={{ disableRipple: true }}
    >
      <Typography variant="h5" class="stat-title">
        Cell Stats
      </Typography>
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
            {`${Math.round(
              dataOverlay.cell_stats.surface_area,
            ).toLocaleString()}nm`}
            <sup>2</sup>
          </p>
        )}
      </Box>
    </AccordionDetails>
  </Accordion>
);

const contactPrimaryNeuron = (contact) => (contact.split('by')[0]);

const ContactStats = ({ dataOverlay, classname }) => (
  <Accordion>
    <AccordionSummary
      expandIcon={
        <img src={CHEVRON} width="auto" height="auto" alt="CHEVRON" />
      }
      IconButtonProps={{ disableRipple: true }}
    >
      <Typography variant="h5" class="stat-title">
        Contact Stats
      </Typography>
    </AccordionSummary>
    <AccordionDetails>
      <Box className="data-overlay-body">
        {dataOverlay.ranking?.cell_rank && dataOverlay.ranking?.cell_total && (
          <p>
            <strong>
              Contact Rank (per&nbsp;
              {contactPrimaryNeuron(dataOverlay.uid)}
              &nbsp;cell):
              {' '}
            </strong>
            {`${dataOverlay.ranking.cell_rank} of ${dataOverlay.ranking.cell_total}`}
            <HTMLTooltip
              className={classname}
              title={(
                <Typography color="inherit">
                  Cell Rank compares the summed surface area of contacts (&ldquo;patches&rdquo;)
                  between these two neurons relative to all other contact
                  relationships across the surface of
                  {' '}
                  <strong>{contactPrimaryNeuron(dataOverlay.uid)}</strong>
                  . A rank of 1 means
                  this neuron pair shares the largest contact area.
                </Typography>
              )}
            >
              <HelpOutlineIcon />
            </HTMLTooltip>
          </p>
        )}
        {dataOverlay.ranking?.brain_rank && dataOverlay.ranking?.brain_total && (
          <p>
            <strong>Contact Rank (per nerve ring): </strong>
            {`${dataOverlay.ranking.brain_rank} of ${dataOverlay.ranking.brain_total}`}
            <HTMLTooltip
              className={classname}
              title={(
                <Typography color="inherit">
                  Nerve Ring Rank compares the summed surface area of contacts
                  (&ldquo;patches&rdquo;) between these two neurons relative to all other contact
                  relationships for the whole nerve ring. A rank of 1 means
                  this neuron pair shares the largest contact area.
                </Typography>
              )}
            >
              <HelpOutlineIcon />
            </HTMLTooltip>
          </p>
        )}
        {dataOverlay.patch_stats.patch_surface_area ? (
          <p>
            <strong>Total area: </strong>
            {`${Math.round(
              dataOverlay.patch_stats.patch_surface_area,
            ).toLocaleString()}nm`}
            <sup>2</sup>
            <HTMLTooltip
              className={classname}
              title={(
                <>
                  <Typography color="inherit">
                    Total Area = The summed surface area of all
                    contact patches between these two cells.
                  </Typography>
                </>
              )}
            >
              <HelpOutlineIcon />
            </HTMLTooltip>
          </p>
        ) : (
          <p>Below threshold</p>
        )}
        {(dataOverlay.patch_stats.patch_surface_area && dataOverlay.ranking.cell_sa_aggregate) ? (
          <p>
            <strong>
              Contact % (per&nbsp;
              {contactPrimaryNeuron(dataOverlay.uid)}
              &nbsp;cell):
              {' '}
            </strong>
            {dataOverlay.ranking.cell_sa_aggregate && (
              <>
                {`${(
                  (dataOverlay.patch_stats.patch_surface_area
                    / dataOverlay.ranking.cell_sa_aggregate)
                  * 100
                ).toFixed(6)}%`}
              </>
            )}
            <HTMLTooltip
              className={classname}
              title={(
                <>
                  <Typography color="inherit">
                    Percent of cell surface area is the contact area of this cell pair to
                    the surface area of
                    {' '}
                    <strong>{contactPrimaryNeuron(dataOverlay.uid)}</strong>
                    .
                  </Typography>
                </>
              )}
            >
              <HelpOutlineIcon />
            </HTMLTooltip>
          </p>
        ) : (
          <p>Below threshold</p>
        )}
        {(dataOverlay.patch_stats.patch_surface_area && dataOverlay.ranking.brain_sa_aggregate) ? (
          <p>
            <strong>Contact % (per nerve ring): </strong>
            {`${(
              (dataOverlay.patch_stats.patch_surface_area
                / dataOverlay.ranking.brain_sa_aggregate)
              * 100
            ).toFixed(6)}%`}
            <HTMLTooltip
              className={classname}
              title={(
                <>
                  <Typography color="inherit">
                    Percent of nerve ring surface area is the contact
                    area of this cell pair to the surface area of the nerve ring.
                  </Typography>
                </>
              )}
            >
              <HelpOutlineIcon />
            </HTMLTooltip>
          </p>
        ) : (
          <p>Below threshold</p>
        )}
      </Box>
    </AccordionDetails>
  </Accordion>
);

const synapseItems = (synapses, classname) => synapses.map((synapse) => (
  <p>
    <strong
      dangerouslySetInnerHTML={{
        __html: formatSynapseUID(synapse.name),
      }}
    />
    :&nbsp;
    {synapse.count}
    <HTMLTooltip
      className={classname}
      title={(
        <>
          <Typography color="inherit">
            Number of synapse similarly configured on this cell.
          </Typography>
        </>
      )}
    >
      <HelpOutlineIcon />
    </HTMLTooltip>
  </p>
));

const SynapseStats = ({ dataOverlay, classname }) => (
  <Accordion>
    <AccordionSummary
      expandIcon={
        <img src={CHEVRON} width="auto" height="auto" alt="CHEVRON" />
      }
      IconButtonProps={{ disableRipple: true }}
    >
      <Typography variant="h5" class="stat-title">
        Synapse Stats
      </Typography>
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
                  (dataOverlay.synapse_stats.total_type_count
                    / dataOverlay.synapse_stats.total_cell_synapse_count)
                    * 100,
                )}%`}
                )
              </>
            )}
          </p>
        )}
        {synapseItems(dataOverlay.synapse_stats.connections, classname)}
      </Box>
    </AccordionDetails>
  </Accordion>
);

const dataOverlayAccordion = (dataOverlay, patchClass) => {
  if (!dataOverlay) {
    return <></>;
  }

  if (dataOverlay?.instanceType === 'neuron') {
    return <CellStats dataOverlay={dataOverlay} />;
  }

  if (dataOverlay?.instanceType === 'contact') {
    return (
      <>
        <ContactStats dataOverlay={dataOverlay} classname={patchClass} />
      </>
    );
  }

  if (dataOverlay?.instanceType === 'synapse') {
    return (
      <>
        <SynapseStats dataOverlay={dataOverlay} classname={patchClass} />
      </>
    );
  }

  return <></>;
};

const dataOverlayTitle = (dataOverlay) => {
  let title = '';
  let parts = [];

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
      title = `<strong>${title}</strong>`;
      break;
    case 'contact':
      title = dataOverlay.uid;
      parts = title.split('by');
      title = `<strong>${parts[0]}</strong> by ${parts[1]}`;
      break;
    case 'synapse':
      title = formatSynapseUID(dataOverlay.uid);
      title = `<strong>${title}</strong>`;
      break;
    default:
      break;
  }

  return title;
};

// const sumSynapses = (synapses) => synapses.reduce((acc, curr) => acc + curr.count, 0);

const DataOverlay = () => {
  const classes = useStyles();
  const storeData = useSelector((state) => state.dataOverlay);
  const { dataOverlay } = storeData;

  return dataOverlay?.uid ? (
    <Box className={classes.root}>
      <Box className="data-overlay">
        <Box className="data-overlay-header">
          <Typography component="h3" className="data-overlay-title">
            <span
              dangerouslySetInnerHTML={{
                __html: dataOverlayTitle(dataOverlay),
              }}
            />
          </Typography>
          <Button
            onClick={() => resetDataOverlay()}
            fontSize="large"
            className="data-overlay-icon"
          >
            <CloseIcon />
          </Button>
        </Box>
        <Divider />
        <Box className="data-overlay-body">
          {dataOverlayAccordion(dataOverlay, classes.tooltip)}
          {/* {dataOverlay.synapses && dataOverlay.synapses.length > 0 && (
            <p>
              <strong>Synapse count: </strong>
              {`${sumSynapses(dataOverlay.synapses)}`}
              &nbsp;(
              {dataOverlay.total_nr_synapses}
              &nbsp;whole nerve ring )
            </p>
          )} */}
        </Box>
      </Box>
    </Box>
  ) : (
    <></>
  );
};

export default DataOverlay;
