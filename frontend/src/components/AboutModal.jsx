import React from 'react';
import Modal from '@material-ui/core/Modal';
import {
  Box,
  Typography,
  IconButton,
} from '@material-ui/core';
import CLOSE from '../images/close.svg';
import { NEUROSCAN_ABOUT, PROMOTERDB_ABOUT } from '../utilities/constants';

const AboutModal = (props) => {
  const { open, handleClose, title } = props;
  let aboutText = '';
  if (title === 'NeuroSCAN') {
    aboutText = NEUROSCAN_ABOUT;
  } else {
    aboutText = PROMOTERDB_ABOUT;
  }
  return (
    <Modal
      open={open}
      className="primary-modal"
      onClose={handleClose}
    >
      <Box className="modal-dialog">
        <Box className="modal-header">
          <Typography>{`About ${title}`}</Typography>
          <IconButton
            color="inherit"
            onClick={handleClose}
            disableFocusRipple
            disableRipple
          >
            <img src={CLOSE} alt="Close" />
          </IconButton>
        </Box>
        <Box className="modal-body">
          {aboutText.map((value) => (
            <div
              dangerouslySetInnerHTML={{
                __html: value,
              }}
            />
          ))}
          <h3>Website design and development:</h3>
          <Typography>
            {`${title} was initially developed by MetaCell and is currently maintained by Bilte Co.`}
          </Typography>
        </Box>
        <Box className="modal-footer">
          {/* <Button variant="outlined" onClick={handleClose}>
            Legal Notice
          </Button>
          <Button disableElevation color="primary" variant="contained" onClick={handleClose}>
            Share
          </Button> */}
        </Box>
      </Box>
    </Modal>
  );
};

export default AboutModal;
