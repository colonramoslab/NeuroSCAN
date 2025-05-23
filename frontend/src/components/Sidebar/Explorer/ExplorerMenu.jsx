import React from 'react';
import {
  Divider,
  List,
  ListItem,
  ListItemText,
  Popover,
  Typography,
} from '@material-ui/core';
import { useSelector } from 'react-redux';
import MenuGroups from './MenuGroups';
import MENU from '../../../images/menu-chevron.svg';
import {
  deleteSelectedInstances,
  handleSelect,
  hideSelectedInstances,
  showSelectedInstances,
} from '../../../services/instanceHelpers';

const ExplorerMenu = ({
  anchorEl, handleMenuClose, open, viewerId, instance, ...other
}) => {
  const widgets = useSelector((state) => state.widgets);

  const handleSelectClick = () => {
    handleSelect(viewerId, instance, widgets);
    handleMenuClose();
  };

  const handleDeleteClick = async () => {
    if (viewerId) {
      await handleSelect(viewerId, instance, widgets);
      deleteSelectedInstances(viewerId, [instance.uid]);
    }
    handleMenuClose();
  };

  const handleHideClick = () => {
    if (viewerId) {
      const { instances } = widgets[viewerId].config;
      hideSelectedInstances(viewerId, instances, [instance.uid]);
    }
    handleMenuClose();
  };

  const handleShowClick = () => {
    if (viewerId) {
      const { instances } = widgets[viewerId].config;
      showSelectedInstances(viewerId, instances, [instance.uid]);
    }
    handleMenuClose();
  };

  const menuId = 'explorer-menu-option';
  const showHideButton = instance.hidden ? (
    <ListItem
      onClick={handleShowClick}
      role="button"
      button
    >
      <ListItemText>
        <Typography>Show</Typography>
      </ListItemText>
    </ListItem>
  ) : (
    <ListItem
      onClick={handleHideClick}
      role="button"
      button
    >
      <ListItemText>
        <Typography>Hide</Typography>
      </ListItemText>
    </ListItem>
  );
  return (
    <Popover
      className="custom-popover dark right no-pin"
      id={menuId}
      open={open}
      anchorEl={anchorEl}
      onClose={handleMenuClose}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'left',
      }}
      transformOrigin={{
        vertical: 'top',
        horizontal: 'left',
      }}
    >
      <List>
        <ListItem
          onClick={handleSelectClick}
          role="button"
          button
        >
          <ListItemText>
            <Typography>Select</Typography>
          </ListItemText>
        </ListItem>

        <Divider />

        <ListItem
          role="button"
          button
        >
          <ListItemText>
            <Typography>
              Add to group
              <img src={MENU} alt="Menu" />
            </Typography>
          </ListItemText>
          <MenuGroups
            instance={instance}
            viewerId={viewerId}
            {...other}
          />
        </ListItem>

        <Divider />

        {showHideButton}

        <ListItem
          onClick={handleDeleteClick}
          role="button"
          disableGutters
          button
        >
          <ListItemText>
            <Typography>Delete</Typography>
          </ListItemText>
        </ListItem>
      </List>
    </Popover>
  );
};

export default ExplorerMenu;
