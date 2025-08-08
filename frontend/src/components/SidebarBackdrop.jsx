import React from 'react';
import { makeStyles } from '@material-ui/core';

const useStyles = makeStyles(() => ({
  root: {
    position: 'fixed',
    top: 0,
    left: 0,
    width: '100%',
    height: '100%',
    backgroundColor: 'rgba(0, 0, 0, 0.5)',
    zIndex: 1000,

    '@media screen and (min-width: 600px)': {
      display: 'none',
    },
  },
}));

const SidebarBackdrop = ({ shrink, setShrink }) => {
  const classes = useStyles();
  const handleClick = (event) => {
    setShrink(true);
    event.stopPropagation();
  };

  const handleKeyPress = () => (event) => {
    if (event.key === 'Escape') {
      setShrink(true);
    }
  };

  return (
    <>
      {/* eslint-disable-next-line jsx-a11y/no-static-element-interactions */}
      {!shrink && <div className={classes.root} onClick={handleClick} onKeyUp={handleKeyPress} />}
    </>
  );
};

export default SidebarBackdrop;
