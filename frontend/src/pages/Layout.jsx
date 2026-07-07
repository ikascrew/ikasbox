import React from "react";

import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import Drawer from '@mui/material/Drawer';

import List from '@mui/material/List';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';

import Pages from "./Pages.jsx";

import Dialog from "./LayoutDialog.jsx";

const drawerWidth = 220;

var inst;

class Layout extends React.Component {

  constructor(props) {
    super(props);
    inst = this;

    this.globalDialog = React.createRef();
    this.showDialog = this.showDialog.bind(this);
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
    } else if ( url.indexOf("/contents") !== -1 ) {
      return "contents";
    }
    return "index";
  }

  showDialog(title,msg,type) {
    return this.globalDialog.current.show(title,msg,type);
  }

  render() {

    var selected = this.getLocation();

    return (
      <Box sx={{ display: 'flex' }}>

      <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            <Box component="a" href="/" sx={{ color: 'inherit', textDecoration: 'none' }}>ikasbox</Box>
          </Typography>
        </Toolbar>
      </AppBar>

      <Dialog ref={this.globalDialog} />

      <Drawer
        variant="permanent"
        sx={{
          width: drawerWidth,
          flexShrink: 0,
          [`& .MuiDrawer-paper`]: { width: drawerWidth, boxSizing: 'border-box' },
        }}
      >
        <Toolbar />
        <List component="nav">

          <ListItemButton
            selected={selected === "category"}
            onClick={(event) => this.handleListItemClick(event, "category")}
          >
            <ListItemIcon></ListItemIcon>
            <ListItemText primary="Category" />
          </ListItemButton>

          <ListItemButton
            selected={selected === "tags"}
            onClick={(event) => this.handleListItemClick(event, "tags")}
          >
            <ListItemIcon></ListItemIcon>
            <ListItemText primary="Tag" />
          </ListItemButton>

          <ListItemButton
            selected={selected === "contents"}
            onClick={(event) => this.handleListItemClick(event, "contents")}
          >
            <ListItemIcon></ListItemIcon>
            <ListItemText primary="Contents" />
          </ListItemButton>

          <ListItemButton
            selected={selected === "groups"}
            onClick={(event) => this.handleListItemClick(event, "groups")}
          >
            <ListItemIcon></ListItemIcon>
            <ListItemText primary="Group" />
          </ListItemButton>

          <ListItemButton
            selected={selected === "projects"}
            onClick={(event) => this.handleListItemClick(event, "projects")}
          >
            <ListItemIcon></ListItemIcon>
            <ListItemText primary="Project" />
          </ListItemButton>

        </List>
      </Drawer>

      <Box component="main" sx={{ flexGrow: 1, p: 3 }}>
        <Toolbar />
        <Pages />
      </Box>

      </Box>
    );
  }
}

export function Alert(title,msg) {
  return inst.showDialog(title,msg,"alert");
}

export function Confirm(title,msg) {
  return inst.showDialog(title,msg,"confirm");
}

export default Layout;
