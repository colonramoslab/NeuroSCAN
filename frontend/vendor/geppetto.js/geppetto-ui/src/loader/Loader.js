import React, { Component, Fragment } from 'react';
import { withStyles } from '@material-ui/core';
import PropTypes from 'prop-types';
import Backdrop from '@material-ui/core/Backdrop';
import CircularProgress from '@material-ui/core/CircularProgress';
import LinearProgress from '@material-ui/core/LinearProgress';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';

const styles = theme => ({
  backdrop: { zIndex: theme.zIndex.drawer + 1, },
  root: {
    position: 'absolute',
    flex: '0 0 100%',
    alignSelf: 'stretch',
  },
});

class Loader extends Component {
  constructor (props) {
    super(props);
    this.state = { messageIndex: 0, };
  }

  render () {
    const {
      classes,
      active,
      fullscreen,
      handleClose,
      messages,
      elapsed,
      backgroundStyle,
      children,
      className,
      messagesInterval
    } = this.props;
    
    if (messagesInterval && messages.length) {
      setTimeout(() => {
        const { messageIndex } = this.state;
        const { messages, active } = this.props;
        if (messages.length && active) {
          this.setState({ messageIndex: (messageIndex + 1) % messages.length });
        }
        
      }, messagesInterval);
    }

    const message = messages.length > 0 ? messages[this.state.messageIndex % messages.length] : '';
    const progress = elapsed ? (
      <LinearProgress
        variant="determinate"
        value={elapsed * 100}
        style={{ width: '200px' }}
      />
    ) : (
      <CircularProgress color="inherit" />
    );

    const content = children ? (
      children
    ) : (
      <Grid container spacing={1}>
        <Grid container item spacing={3} justify="center">
          <Grid item>{progress}</Grid>
        </Grid>
        <Grid container item spacing={3} justify="center">
          <Grid item>
            <Typography
              className={className}
              display="block"
              variant="subtitle1"
              gutterBottom
            >
              {message}
            </Typography>
          </Grid>
        </Grid>
      </Grid>
    );

    const backdrop = fullscreen ? (
      <Fragment>
        <Backdrop
          className={classes.backdrop}
          open={active}
          onClick={handleClose}
          style={backgroundStyle}
        >
          {content}
        </Backdrop>
      </Fragment>
    ) : (
      <Fragment>
        <Backdrop
          open={active}
          onClick={handleClose}
          style={backgroundStyle}
          className={classes.backdrop}
          classes={{ root: classes.root, }} // class name, e.g. `classes-nesting-root-x` 
        >
          {content}
        </Backdrop>
      </Fragment>
    );

    return backdrop;
  }
}

Loader.defaultProps = {
  active: true,
  fullscreen: true,
  messages: [],
  messagesInterval: 10000,
  elapsed: null,
  backgroundStyle: { backgroundColor: 'rgba(255,142,0,0.1)' },
  handleClose: () => {},
};

Loader.propTypes = {
  /**
   * Flag to show/hide the Loader
   */
  active: PropTypes.bool,
  /**
   * Flag to show/hide the Loader in fullscreen
   */
  fullscreen: PropTypes.bool,
  /**
   * Function to handle the close of the Loader
   */
  handleClose: PropTypes.func,
  /**
   * Array of Custom messages to display
   */
  messages: PropTypes.array,
  /**
   * Number of milliseconds between custom messages
   */
  messagesInterval: PropTypes.number,
  /**
   * Number of the progress value to show in linear determinate (in percentage)
   */
  elapsed: PropTypes.number,
  /**
   * Style to be applied to the Loader background
   */
  backgroundStyle: PropTypes.shape({
    /**
     * Loader's background color. Defaults to rgba(255,142,0,0.1)
     */
    backgroundColor: PropTypes.string,
  }),
};

export default withStyles(styles)(Loader);
