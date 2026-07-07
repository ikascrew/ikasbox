import React from "react"

import {withParams} from "../Pages.jsx"

import API from "../../API";
import Util from "../../Util";

import {
  Button, Dialog, DialogTitle, DialogContent, List, ListItemButton, ListItemText
} from '@mui/material';

import FlexTable from "../../components/FlexTable";

class ProjectGroups extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      groups : new Map(),
      addOpen : false
    }

    this.columns = [
      { id: 'id', label: 'ID', minWidth: 20 },
      { id: 'name', label: 'Name', minWidth: 100,
        format: (val, row) => this.contentLink(val, row) },
      { id: 'created_at', label: 'Created At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
      { id: 'updated_at', label: 'Updated At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
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
        groups : groups
      });

      this.table.current.set(result.groups);

    }).catch((err) => {
      console.log(err)
    });
  }

  contentLink(val, row) {
    return <a href={"/groups/contents/" + row["id"]}>{val}</a>;
  }

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

  render() {

    var groups = this.state.groups;

    return (<>
      <Button variant="contained" onClick={this.handleOpenAdd}>+</Button>

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
    </>);
  }
}

export default withParams(ProjectGroups);
