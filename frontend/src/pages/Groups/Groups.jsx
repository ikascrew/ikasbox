import React, { createRef } from "react"

import API from "../../API.js"
import Util from "../../Util.js"
import Paging from "../Paging.js";

import {
  Box,
  Button
} from '@mui/material';

import GroupRegisterDialog from "./GroupRegisterDialog"
import { Confirm, Alert } from "../Layout.jsx";
import FlexTable from "../../components/FlexTable.jsx";

class Groups extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      paging:Paging.create()
    }

    this.handleChangePage = this.handleChangePage.bind(this);
    this.handleChangeRowsPerPage = this.handleChangeRowsPerPage.bind(this);
    this.handleOpenRegister = this.handleOpenRegister.bind(this);

    this.registerDialog = React.createRef();

    this.columns = [
      { id: 'name', label: 'Name', minWidth: 100,
        format: (val, row) => this.nameCell(val, row) },
      { id: 'created_at', label: 'Created At', minWidth: 160, width: 160,
        format: (value) => Util.formatDate(value) },
      { id: 'updated_at', label: 'Updated At', minWidth: 160, width: 160,
        format: (value) => Util.formatDate(value) },
      { id: 'delete', label: '', minWidth: 100, width: 100,
        format: (val, row) => this.createDeleteButton(val, row) },
    ];

    this.table = React.createRef();
  }

  componentDidMount() {
    var paging = this.state.paging;
    this.view(paging);
  }

  /**
   * 一覧の更新
   * @param {*} paging 
   */
  view(paging) {

    var args = {
      paging: paging
    }

    API.post("/api/v1/groups/view", args).then((res) => {
      var result = res.data;
      this.table.current.set(result.groups, result.paging);
      this.setState({
        paging:result.paging
      });
    }).catch((err) => {
      console.log(err)
    });
  }

  handleOpenRegister() {
    this.registerDialog.current.open();
  }

  render() {

    return (<>

      <GroupRegisterDialog ref={this.registerDialog} onCommit={() => this.view(this.state.paging)} />
      <Button variant="contained" onClick={this.handleOpenRegister}>Register</Button>
      <FlexTable columns={this.columns} paging={this.state.paging} ref={this.table}
        onPageChange={this.handleChangePage} onRowsPerPageChange={this.handleChangeRowsPerPage}/>

    </>);
  }

  handleChangePage = (event, newPage) => {
    var paging = this.state.paging;
    paging.current = newPage;
    this.view(paging);
  };

  handleChangeRowsPerPage = (event) => {
    var paging = this.state.paging;
    paging.current = 1;
    paging.limit = event.target.value;
    this.view(paging);
  }


  contentLink(val, row) {
    return <a href={"/groups/contents/" + row["id"]}>{val}</a>;
  }

  nameCell(val, row) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        {this.contentLink(val, row)}
        {this.createCheckButton(val, row)}
      </Box>
    );
  }

  handleCheck(val) {

    Confirm("Group Check", "Check for missing files in \"" + val.name + "\"?").then(() => {

      var args = {
        groupId: val.id
      }

      API.post("/api/v1/groups/check", args).then((res) => {
        var result = res.data;
        if (result.missing.length <= 0) {
          Alert("Group Check", "All " + result.total + " content file(s) exist.");
          return;
        }

        var lines = result.missing.map((c) => c.id + ": " + c.name + " (" + c.path + ")");
        Alert("Group Check",
          "Missing " + result.missing.length + "/" + result.total + " file(s):\n" + lines.join("\n"));

      }).catch((err) => {
        console.log(err);
      });

    }).catch((e) => {
      console.log(e);
    });
  }

  createCheckButton(val, row) {
    return (
      <Button color="warning" variant="contained" onClick={() => this.handleCheck(row)}>
        Check
      </Button>
    );
  }

  handleDelete(val) {

    Confirm("Group Delete", "Delete group \"" + val.name + "\" and all its contents?").then(() => {

      var args = {
        groupId: val.id
      }

      API.delete("/api/v1/groups/delete", { data: args }).then(() => {
        this.view(this.state.paging);
      }).catch((err) => {
        console.log(err);
      });

    }).catch((e) => {
      console.log(e);
    });
  }

  createDeleteButton(val, row) {
    return (
      <Button color="error" variant="contained" onClick={() => this.handleDelete(row)}>
        Delete
      </Button>
    );
  }
}

export default Groups;