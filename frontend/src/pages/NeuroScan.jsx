import React, { useEffect, useState, useCallback } from 'react';
import { useSelector, useStore, useDispatch } from 'react-redux';
import { Box, CircularProgress, makeStyles } from '@material-ui/core';
import { getLayoutManagerInstance } from '@metacell/geppetto-meta-client/common/layout/LayoutManager';
import LeftSidebar from '../components/LeftSidebar';
import SidebarBackdrop from '../components/SidebarBackdrop';
import DataOverlay from '../components/DataOverlay/DataOverlay';
import Header from '../components/Header';
import { VIEWS, VIEWERS } from '../utilities/constants';
import ViewerPlaceholder from '../components/ViewerPlaceholder';
import { addInstances } from '../redux/actions/widget';
import { mapLocalGltfToInstance } from '../services/instanceHelpers';

const useStyles = makeStyles((theme) => ({
  root: {
    [theme.breakpoints.down('xs')]: {
      '& .primary-structure': {
        paddingTop: '2.5rem',
        display: 'block !important',
      },
    },
  },
  layoutContainer: {
    position: 'relative',
    width: '100%',
    height: '100%',
  },
  left: {
    flexShrink: 0,
  },
  right: {
    flexGrow: 1,
    [theme.breakpoints.down('xs')]: {
      height: '100%',
    },
  },
  dropOverlay: {
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    backgroundColor: 'rgba(44, 44, 44, 0.9)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
    border: '3px dashed #a0c020',
    borderRadius: '8px',
    pointerEvents: 'none',
  },
  dropText: {
    color: '#a0c020',
    fontSize: '1.5rem',
    fontWeight: 'bold',
  },
}));

const readFileAsBase64 = (file) => new Promise((resolve, reject) => {
  const reader = new FileReader();
  reader.onload = () => {
    const dataUrl = String(reader.result || '');
    const comma = dataUrl.indexOf(',');
    resolve(comma >= 0 ? dataUrl.slice(comma + 1) : dataUrl);
  };
  reader.onerror = reject;
  reader.readAsDataURL(file);
});

export default function NeuroScan() {
  const classes = useStyles();
  const store = useStore();
  const dispatch = useDispatch();
  const [LayoutComponent, setLayoutManager] = useState(undefined);
  const [shrinkSidebar, setShrinkSidebar] = React.useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const widgets = useSelector((state) => state.widgets);
  const viewerCount = Object.keys(widgets).length;

  const handleToggle = () => {
    setShrinkSidebar(!shrinkSidebar);
  };

  const handleDragEnter = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.dataTransfer.types.includes('Files')) {
      setIsDragging(true);
    }
  }, []);

  const handleDragOver = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
  }, []);

  const handleDragLeave = useCallback((e) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.currentTarget.contains(e.relatedTarget)) {
      return;
    }
    setIsDragging(false);
  }, []);

  const handleDrop = useCallback(async (e) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);

    const { files } = e.dataTransfer;
    if (!files || files.length === 0) return;

    const gltfFiles = Array.from(files).filter((file) => {
      const name = file.name.toLowerCase();
      return name.endsWith('.gltf') || name.endsWith('.glb');
    });

    if (gltfFiles.length === 0) return;

    const existingViewerId = viewerCount > 0 ? Object.keys(widgets)[0] : null;

    const instances = await Promise.all(
      gltfFiles.map(async (file) => {
        const base64 = await readFileAsBase64(file);
        return mapLocalGltfToInstance({ fileName: file.name, base64 });
      }),
    );

    dispatch(addInstances(existingViewerId, instances, VIEWERS.InstanceViewer));
  }, [dispatch, viewerCount, widgets]);

  useEffect(() => {
    const preventDefaultDrop = (e) => {
      e.preventDefault();
    };
    window.addEventListener('dragover', preventDefaultDrop);
    window.addEventListener('drop', preventDefaultDrop);
    return () => {
      window.removeEventListener('dragover', preventDefaultDrop);
      window.removeEventListener('drop', preventDefaultDrop);
    };
  }, []);

  useEffect(() => {
    if (LayoutComponent === undefined) {
      const myManager = getLayoutManagerInstance();
      if (myManager) {
        setLayoutManager(myManager.getComponent());
      }
    }
  }, [store]);

  let componentToRender = <CircularProgress />;
  if (LayoutComponent !== undefined) {
    if (viewerCount === 0) {
      componentToRender = <ViewerPlaceholder />;
    } else {
      componentToRender = <LayoutComponent />;
    }
  }

  return (
    <Box className={classes.root}>
      <Box className="primary-structure" display="flex">
        <Box className={classes.left}>
          <Header shrink={shrinkSidebar} toggleSidebar={handleToggle} view={VIEWS?.neuroScan} />
          <LeftSidebar shrink={shrinkSidebar} />
          <SidebarBackdrop shrink={shrinkSidebar} setShrink={setShrinkSidebar} />
        </Box>
        <Box className={classes.right}>
          <Box
            className={`primary-structure_content ${classes.layoutContainer} ${viewerCount > 0 ? 'padding' : ''}`}
            onDragEnter={handleDragEnter}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
          >
            {isDragging && (
              <Box className={classes.dropOverlay}>
                <span className={classes.dropText}>Drop .gltf/.glb file to load</span>
              </Box>
            )}
            {componentToRender}
          </Box>
        </Box>
        <DataOverlay />
      </Box>
    </Box>
  );
}
