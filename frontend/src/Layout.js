import React from "react";

import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';

import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Grid from '@mui/material/Grid';

import "./css/Layout.css"

class Layout extends React.Component {
  render() {
    var selectedIndex = 0;
    return (
      <>
      <AppBar position="static">
        <Toolbar>

          <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}
          />

          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            ikasbox
          </Typography>

        </Toolbar>
      </AppBar>

      <Grid container spacing={1} className="MainGrid">

        <Grid item xs={3} className="ListGrid">
          <List component="nav" aria-label="main mailbox folders">
            <ListItemButton
              selected={selectedIndex === 0}
              onClick={(event) => handleListItemClick(event, 0)}
            >
              <ListItemIcon>
              </ListItemIcon>
              <ListItemText primary="Group" />
            </ListItemButton>
  
            <ListItemButton
              selected={selectedIndex === 1}
              onClick={(event) => handleListItemClick(event, 1)}
            >
              <ListItemIcon>
              </ListItemIcon>
              <ListItemText primary="Project" />
            </ListItemButton>
          </List> 
        </Grid>

        <Grid item xs={9} className="ContentGrid"> 

Router

        </Grid>
      </Grid>
    </>);
  }
}

export default Layout;

