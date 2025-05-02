import * as React from 'react';
import { withStyles } from '@material-ui/core/styles';
import Tooltip from '@material-ui/core/Tooltip';

const HtmlTooltip = withStyles((theme) => ({
  tooltip: {
    backgroundColor: '#f5f5f9',
    color: 'rgba(0, 0, 0, 1)',
    maxWidth: 220,
    fontSize: theme.typography.pxToRem(16),
    border: '1px solid #dadde9',
  },
}))(Tooltip);

export default HtmlTooltip;
