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

export const downloadVideo = (url, filename) => {
  // eslint-disable-next-line no-undef
  const a = document.createElement('a');
  a.style.display = 'none';
  a.href = url;
  a.download = filename;
  // eslint-disable-next-line no-undef
  document.body.appendChild(a);
  a.click();
  setTimeout(() => {
    // eslint-disable-next-line no-undef
    document.body.removeChild(a);
    window.URL.revokeObjectURL(url);
  }, 100);
};

const RecordControlModal = (props) => {
  const {
    open, handleClose, videoBlob, widgetName,
  } = props;
  const [deleteOption, setDeleteOption] = useState(false);
  const [processing, setProcessing] = useState(false);
  const [error, setError] = useState(null);
  const [showDownload, setShowDownload] = useState(true);
  const downloadRecording = async () => {
    setProcessing(true);
    setShowDownload(false);
    // const mp4 = Buffer.from(await webmToMp4(Buffer.from(await videoBlob.arrayBuffer())));
    // const videoBuffer = await videoBlob.arrayBuffer();
    const videoMeta = await webmToMp4(videoBlob);
    const videoUUID = videoMeta.id;

    // poll the server to check if the mp4 is ready at /videos/status/:id every 2 seconds

    const checkStatus = async (retries = 30) => {
      if (retries === 0) {
        console.error('Max retries reached');
        return null;
      }

      // wait for 2 seconds
      await new Promise((resolve) => setTimeout(resolve, 2000));
      const response = await fetch(`${process.env.REACT_APP_BACKEND_URL}videos/status/${videoUUID}`);
      // if the response.status is any 200
      if (!response.ok) {
        console.error('Error fetching video status:', response.status);
        setError('Error fetching video status');
        return null;
      }
      // if the response is any 200

      if (response.status >= 200 && response.status < 300) {
        const data = await response.json();
        if (data.status === 'completed') {
          return data;
        }

        if (data.status === 'failed') {
          console.error('Video conversion failed');
          setError('Video conversion failed');
          return null;
        }

        return checkStatus(retries - 1);
      }

      console.error('Error checking status:', response.status);
      return null;
    };

    const status = await checkStatus();

    if (!status) {
      setShowDownload(true);
      return;
    }

    if (status.status === 'completed') {
      downloadVideo(`${process.env.REACT_APP_BACKEND_URL}videos/download/${videoUUID}.mp4`, `${videoUUID}.mp4`);
      setShowDownload(true);
      handleClose();
    } else {
      setError('Video conversion failed');
    }
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
          { processing
            && (
              <Typography>
                Processing...
              </Typography>
            ) }
          { Error
            && (
              <Typography>
                { error }
              </Typography>
            ) }
        </Box>
      </Box>
    </Modal>
  );
};

export default RecordControlModal;
