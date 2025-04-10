import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Box, Button, Typography } from '@material-ui/core';
import CircularLoader from '../Common/Loader';
import AddToViewerMenu from './AddToViewerMenu';
import SearchResult from '../Common/SearchResult';
import NEURON from '../../images/neuron.svg';
import CONTACTS from '../../images/contacts.svg';
import SYNAPSES from '../../images/synapses.svg';
import { addInstances } from '../../redux/actions/widget';
import { mapToInstance } from '../../services/instanceHelpers';
import { VIEWERS } from '../../utilities/constants';
import * as search from '../../redux/actions/search';

const list = [
  {
    title: 'Cells',
    resultItem: 'neurons',
    image: NEURON,
  },
  {
    title: 'Contacts',
    resultItem: 'contacts',
    image: CONTACTS,
  },
  {
    title: 'Synapses',
    resultItem: 'synapses',
    image: SYNAPSES,
  },
];

const initialSelectedItems = {
  neurons: [],
  contacts: [],
  synapses: [],
  scale: [],
};
const Results = ({ timePoint }) => {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [currentItem, setCurrentItem] = React.useState(null);
  const [selectedItems, setSelectedItems] = useState(initialSelectedItems);
  const dispatch = useDispatch();

  const searchesCount = useSelector((state) => state.search.searchesCount);
  const scaleItems = useSelector((state) => state.search.results.scale);
  const handleClick = (event, selectedItem) => {
    setCurrentItem(selectedItem);
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setCurrentItem(null);
    setAnchorEl(null);
  };

  const setScaleItem = () => {
    // dispatch(search.getAll({ entity: 'scale' }));
    dispatch(search.loadMore({ entity: 'scale' }));
    if (scaleItems.items && scaleItems.items.length > 0) {
      setSelectedItems({
        ...selectedItems,
        scale: [scaleItems.items[0]],
      });
    }
  };

  const handleAddToViewer = async (viewerId = null) => {
    setScaleItem();
    if (currentItem) {
      const instances = [mapToInstance(currentItem)];
      dispatch(addInstances(viewerId, instances, VIEWERS.InstanceViewer));
    } else if (Object.values(selectedItems).some((array) => array.length > 0)) {
      const itemsArray = Object.values(selectedItems).flat();
      const instances = itemsArray.map((item) => mapToInstance(item));
      dispatch(addInstances(viewerId, instances, VIEWERS.InstanceViewer));
      setSelectedItems(initialSelectedItems);
      Object.keys(selectedItems).forEach((key) => dispatch(search.deselectAll({ entity: key })));
    }
    handleClose();
  };

  let buttonComponent = null;

  if (searchesCount > 0) {
    buttonComponent = <CircularLoader />;
  } else if (Object.values(selectedItems).some((array) => array.filter((item) => item.instanceType !== 'scale').length > 0)) {
    buttonComponent = (
      <Button
        disableElevation
        aria-haspopup="true"
        color="primary"
        variant="contained"
        onClick={(e) => handleClick(e, null)}
      >
        Add Selected
      </Button>
    );
  }

  useEffect(() => {
    setSelectedItems(initialSelectedItems);
    setScaleItem();
  }, [timePoint]);

  return (
    <Box className="wrap">
      <Typography component="h3">
        Results
        { buttonComponent }
      </Typography>
      {
          list.map((record, index) => (
            <SearchResult
              key={`results-${index}`}
              title={record?.title}
              resultItem={record?.resultItem}
              image={record?.image}
              service={record?.service}
              handleClick={handleClick}
              selectedItems={selectedItems}
              setSelectedItems={setSelectedItems}
              timePoint={timePoint}
            />
          ))
        }
      <AddToViewerMenu
        handleClose={handleClose}
        handleAddToViewer={handleAddToViewer}
        anchorEl={anchorEl}
      />
    </Box>
  );
};

export default Results;
