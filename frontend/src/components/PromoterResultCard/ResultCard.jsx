import React from 'react';
import {
  Box,
  Typography,
  Chip,
  Grid,
} from '@material-ui/core';
import TICK from '../../images/tick.svg';
import ResultTabs from './ResultTabs';
import { filesURL, PROMOTER_MEDIA_TYPES } from '../../utilities/constants';
import PromoterInformation from './PromoterInformation';

const ResultCard = ({
  result,
}) => {
  const promoterBaseUrl = `${filesURL}/promoters/promoter${result.uid.toLowerCase()}`;
  const timeline = [{
    label: 'Timeline',
    src: `${promoterBaseUrl}/timeline.svg`,
  }];
  const model = [{
    label: 'Model',
    src: `${promoterBaseUrl}/model.svg`,
  }];
  const expression = [
    {
      label: '3D Expression',
      src: `${promoterBaseUrl}/3d_expression.mp4`,
      mediaType: PROMOTER_MEDIA_TYPES.video,
    },
  ];
  const promoterVideos = [
    {
      label: 'Promoter',
      src: `${promoterBaseUrl}/promoter.mp4`,
      mediaType: PROMOTER_MEDIA_TYPES.video,
    },
    {
      label: 'Histone Marker',
      src: `${promoterBaseUrl}/histone_marker.mp4`,
      mediaType: PROMOTER_MEDIA_TYPES.video,
    },
  ];
  const cellsByLineaging = result.cellsByLineaging ? result.cellsByLineaging.split(' ').filter((c) => c !== '') : [];
  const otherCells = result.otherCells ? result.otherCells.split(' ').filter((c) => c !== '') : [];
  return (
    <Box className="results-box">
      <Box className="results-box_header">
        <Typography component="h3">
          {result.uid}
        </Typography>
        <Box className="wrap">
          <Box className="tags" display="flex" flexWrap="wrap" justifyContent="flex-end">
            {
              cellsByLineaging.map((cell, index) => <Chip key={`celllineage_${index}`} avatar={<img src={TICK} alt="tick" />} label={cell} className="active" />)
            }
            {
              otherCells.map((cell, index) => <Chip key={`otherCells_${index}`} label={cell} />)
            }
          </Box>
          <Typography>
            <img src={TICK} alt="tick" />
            Cells identified by lineaging
          </Typography>
        </Box>
      </Box>

      <ResultTabs options={timeline} fullWidth />

      <Grid container spacing={1}>
        <Grid item xs={12} sm={4}>
          <ResultTabs options={model} whiteBg />
        </Grid>
        <Grid item xs={12} sm={4}>
          <ResultTabs options={expression} />
        </Grid>
        <Grid item xs={12} sm={4}>
          <ResultTabs options={promoterVideos} />
        </Grid>
      </Grid>

      <PromoterInformation info1={result.information} info2={result.expressionPatterns} />
    </Box>
  );
};

export default ResultCard;
