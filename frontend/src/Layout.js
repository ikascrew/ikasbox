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

import Pages from "./Pages.js";

import "./css/Layout.css";

class Layout extends React.Component {

  constructor(props) {
    super(props);
  }

  handleListItemClick(e,name) {
    location.href="/" + name;
  }

  getLocation() {
    var url = location.href;
    if ( url.indexOf("/groups") !== -1 ) {
      return "groups";
    } else if ( url.indexOf("/projects") !== -1 ) {
      return "projects";
    }

    return "index";
  }

  render() {

    var selected = this.getLocation();

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
            <a href="/">ikasbox</a>
          </Typography>

        </Toolbar>
      </AppBar>

      <Grid container spacing={1} className="MainGrid">

        <Grid item xs={3} className="ListGrid">
          <List component="nav" aria-label="main mailbox folders">


            <ListItemButton
              selected={selected === "category"}
              onClick={(event) => this.handleListItemClick(event, "category")}
            >
              <ListItemIcon>
              </ListItemIcon>
              <ListItemText primary="Category" />
            </ListItemButton>

            <ListItemButton
              selected={selected === "tags"}
              onClick={(event) => this.handleListItemClick(event, "tags")}
            >
              <ListItemIcon>
              </ListItemIcon>
              <ListItemText primary="Tag" />
            </ListItemButton>
  
            <ListItemButton
              selected={selected === "groups"}
              onClick={(event) => this.handleListItemClick(event, "groups")}
            >
              <ListItemIcon>
              </ListItemIcon>
              <ListItemText primary="Group" />
            </ListItemButton>

            <ListItemButton
              selected={selected === "projects"}
              onClick={(event) => this.handleListItemClick(event, "projects")}
            >
              <ListItemIcon>
              </ListItemIcon>
              <ListItemText primary="Project" />
            </ListItemButton>

          </List> 
        </Grid>

        <Grid item xs={9} className="ContentGrid"> 

          <Pages />

        </Grid>
      </Grid>
    </>);
  }
}

export default Layout;

