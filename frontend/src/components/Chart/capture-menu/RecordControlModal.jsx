import React, { useState } from 'react';
import Modal from '@material-ui/core/Modal';
import {
  Box, Typography, IconButton, Button,
} from '@material-ui/core';
import { formatDate } from '@metacell/geppetto-meta-ui/3d-canvas/captureManager/utils';
import CLOSE from '../../../images/close.svg';
import DOWNLOAD from '../../../images/download.svg';
import DELETE from '../../../images/delete.svg';
import DELETE_WHITE from '../../../images/delete-white.svg';

import webmToMp4 from '../../../utilities/webmToMp4';

export const downloadBlob = (blob, filename) => {
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.style.display = 'none';
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  setTimeout(() => {
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  }, 100);
};

const RecordControlModal = (props) => {
  const {
    open, handleClose, videoBlob, widgetName,
  } = props;
  const [deleteOption, setDeleteOption] = useState(false);
  const [showDownload, setShowDownload] = useState(true);
  const downloadRecording = async () => {
    setShowDownload(false);
    // const mp4 = Buffer.from(await webmToMp4(Buffer.from(await videoBlob.arrayBuffer())));
    const videoBuffer = await videoBlob.arrayBuffer();
    const mp4ArrayBuffer = await webmToMp4(videoBuffer);
    // const mp4 = Buffer.from(await webmToMp4(videoBuffer));
    const mp4Blob = new Blob([mp4ArrayBuffer], { type: 'video/mp4' });
    downloadBlob(mp4Blob, `${widgetName}_${formatDate(new Date())}.mp4`);
    setShowDownload(true);
    handleClose();
  };

  const videoSrc = videoBlob ? window.URL.createObjectURL(videoBlob) : null;

  return (
    <Modal
      open={open}
      className="primary-modal"
      onClose={handleClose}
    >
      <Box className="modal-dialog medium">
        <Box className="modal-header">
          <Typography>New screen recording</Typography>
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
          <Box className="video-box">
            {/* eslint-disable-next-line jsx-a11y/media-has-caption */}
            <video src={videoSrc} playsInline controls="controls" className="video-preview" />
          </Box>
        </Box>
        <Box className="modal-footer" justifyContent="space-between">
          { !deleteOption && (
            <Button variant="outlined" onClick={() => setDeleteOption(true)}>
              <img src={DELETE} alt="Delete" />
              Delete
            </Button>
          ) }
          { deleteOption && (
            <Button disableElevation color="secondary" variant="contained" onClick={handleClose}>
              <img src={DELETE_WHITE} alt="Delete" />
              Sure, delete.
            </Button>
          ) }
          { showDownload
            ? (
              <Button disableElevation color="primary" variant="contained" onClick={downloadRecording}>
                <img src={DOWNLOAD} alt="Close" />
                Download
              </Button>
            )
            : null}
        </Box>
      </Box>
    </Modal>
  );
};

export default RecordControlModal;
