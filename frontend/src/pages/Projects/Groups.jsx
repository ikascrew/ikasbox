import React from "react"

import {withParams} from "../Pages.jsx"

import API from "../../API";
import Util from "../../Util";

import {
  TextField, Stack, Button, Dialog, DialogTitle, DialogContent, List, ListItemButton, ListItemText
} from '@mui/material';

import FlexTable from "../../components/FlexTable";
import NameLink from "../../components/NameLink.jsx";
import LoadingButton from "../../components/LoadingButton.jsx";
import { Alert, Confirm } from "../Layout.jsx";

class ProjectGroups extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      groups : new Map(),
      addOpen : false,
      name : ""
    }

    this.columns = [
      { id: 'id', label: 'ID', minWidth: 80, width: 80, align: 'center' },
      { id: 'name', label: 'Name', minWidth: 100,
        format: (val, row) => this.contentLink(val, row) },
      { id: 'created_at', label: 'Created At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
      { id: 'updated_at', label: 'Updated At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
      { id: 'delete', label: '', minWidth: 100, width: 100,
        format: (val, row) => this.createDeleteButton(row) },
    ];

    this.projectId = this.props.params.id;
    this.table = React.createRef();

    this.view = this.view.bind(this);
  }

  componentDidMount() {
    this.view();
  }

  view() {

    var args = {
      projectId : Number(this.projectId)
    }

    API.post("/api/v1/projects/group", args).then((res) => {
      var result = res.data;

      let groups = new Map();
      result.allGroups.forEach( (elm) => {
        groups.set(elm.id,elm.name);
      })

      this.setState({
        groups : groups,
        name : result.project ? result.project.name : ""
      });

      this.table.current.set(result.groups);

    }).catch((err) => {
      console.log(err)
    });
  }

  contentLink(val, row) {
    return <NameLink href={"/groups/contents/" + row["id"]}>{val}</NameLink>;
  }

  handleChangeName = (event) => {
    this.setState({
      name: event.target.value
    });
  };

  handleSaveName = () => {

    var args = {
      projectId: Number(this.projectId),
      name: this.state.name
    }

    return API.patch("/api/v1/projects/rename", args).then(() => {
      Alert("Project", "Name updated.");
    }).catch((err) => {
      console.log(err);
    });
  };

  handleOpenAdd = () => {
    this.setState({ addOpen: true });
  };

  handleCloseAdd = () => {
    this.setState({ addOpen: false });
  };

  handleAddGroup = (groupId) => {

    var args = {
      projectId : Number(this.projectId),
      groupId : Number(groupId)
    }

    API.patch("/api/v1/projects/group/add", args).then((res) => {
      this.handleCloseAdd();
      this.view();
    }).catch((err) => {
      console.log(err)
    });
  }

  handleRemoveGroup(row) {

    Confirm("Remove Group", "Remove group \"" + row.name + "\" from this project?").then(() => {

      var args = {
        projectId: Number(this.projectId),
        groupId: row.id
      }

      API.delete("/api/v1/projects/group/remove", args).then(() => {
        this.view();
      }).catch((err) => {
        console.log(err);
      });

    }).catch((e) => {
      console.log(e);
    });
  }

  createDeleteButton(row) {
    return (
      <Button color="error" variant="contained" onClick={() => this.handleRemoveGroup(row)}>
        Delete
      </Button>
    );
  }

  render() {

    var groups = this.state.groups;

    return (<>

      <Stack direction="row" spacing={2} alignItems="center" sx={{ marginBottom: "10px" }}>
        <TextField
          value={this.state.name}
          onChange={this.handleChangeName}
          autoFocus margin="dense"
          id="name" type="text" label="Name"
          variant="standard"
          InputLabelProps={{ shrink: true }}
          fullWidth
        />
        <LoadingButton onClick={this.handleSaveName}>Save</LoadingButton>
      </Stack>

      <Dialog open={this.state.addOpen} onClose={this.handleCloseAdd}>
        <DialogTitle>Add Group</DialogTitle>
        <DialogContent>
          <List>
{Array.from(groups.entries()).map( ([id, name]) => {
  return (
    <ListItemButton key={id} onClick={() => this.handleAddGroup(id)}>
      <ListItemText primary={name} />
    </ListItemButton>
  );
})}
          </List>
        </DialogContent>
      </Dialog>

      <FlexTable columns={this.columns} ref={this.table} />

      <Button variant="contained" sx={{ marginTop: "10px" }} onClick={this.handleOpenAdd}>+</Button>
    </>);
  }
}

export default withParams(ProjectGroups);
