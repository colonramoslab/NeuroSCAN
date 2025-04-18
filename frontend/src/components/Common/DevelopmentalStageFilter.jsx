import React from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  makeStyles,
  Box,
  Slider,
  Typography,
} from '@material-ui/core';
import DOWN from '../../images/chevron-down.svg';
import vars from '../../styles/constants';
import * as search from '../../redux/actions/search';

const useStyles = makeStyles(() => ({
  root: {
    height: '3rem',
    display: 'flex',
    justifyContent: 'center',
    flexDirection: 'column',
    '&.fade': {
      opacity: 0.3,
      filter: 'grayscale(1)',
      pointerEvents: 'none',
    },
  },
  sliderValue: {
    display: 'flex',
    alignItems: 'center',
    position: 'relative',
    '& p': {
      letterSpacing: '0.005em',
      fontSize: '.6rem',
      lineHeight: '.6rem',
      fontWeight: 'bold',
      color: vars?.captionTextColor,
      textAlign: 'center',
      flexShrink: 0,
      paddingLeft: '2px',
      borderLeft: `1px solid ${vars.sliderBorderColor}`,
      '&:last-child': {
        borderRight: `1px solid ${vars.sliderBorderColor}`,
      },
    },
  },
}));

const sliderMarker = <img width="6" height="4" src={DOWN} alt="DOWN" />;

const DevelopmentalStageFilter = (props) => {
  const { timePoint, setTimePoint } = props;
  const classes = useStyles();
  const dispatch = useDispatch();

  const [sliderVal, setSliderVal] = React.useState(timePoint);
  const devStages = useSelector((state) => state.devStages.neuroSCAN);

  const handleChange = (e, value) => {
    setTimePoint(value);
    dispatch(search.deselectAll());
  };

  const marks = devStages
    .reduce((x, devStage) => (x.concat(devStage.timepoints)), [])
    .filter((item) => item !== undefined)
    .map((mark) => ({
      value: parseInt(mark, 10),
      label: sliderMarker,
    }));

  const min = Math.min(...devStages.map((devStage) => devStage.begin));
  // eslint-disable-next-line max-len
  const max = Math.max(...devStages.map((devStage) => Math.max(devStage.end, devStage.begin)));
  devStages.sort((a, b) => a.order - b.order);
  const stepWidth = (max - min) / 100;

  return devStages.length > 0 && (
    <Box className={classes.root}>
      <Slider
        defaultValue={sliderVal}
        aria-labelledby="developmental-stage-filter"
        step={null}
        valueLabelDisplay="auto"
        marks={marks}
        min={min}
        max={max}
        aria-label="Developmental Stages Filter"
        onChange={handleChange}
      />
      <Box className={classes.sliderValue}>
        {
          devStages.map((stage) => {
            const stageWidth = (Math.max(stage.end, stage.begin) - stage.begin) / stepWidth;
            return stageWidth > 0 && (
              <Typography
                key={stage.id}
                style={{ width: `${stageWidth}%` }}
              >
                {stage.uid}
              </Typography>
            );
          })
        }
      </Box>
    </Box>
  );
};

export default DevelopmentalStageFilter;
