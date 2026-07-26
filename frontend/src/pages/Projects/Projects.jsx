import React from "react"
import Paging from "../Paging";
import API from "../../API";
import Util from "../../Util";

import {
  Box,
  Button
} from '@mui/material';

import FlexTable from "../../components/FlexTable";
import NameLink from "../../components/NameLink.jsx";
import ProjectRegisterDialog from "./ProjectRegisterDialog";
import { Alert, Confirm } from "../Layout.jsx";
import { Link as RouterLink } from "react-router";

class Projects extends React.Component {
  
  constructor(props) {
    super(props);

    this.state = {
      data : [],
      paging : Paging.create(),
      // ika-server -ikasbox 同居モード(v1/server/status が生えている)
      // のときだけ Server ボタン列を表示する
      serverMode : false,
      serverProjectId : 0
    }

    this.columns = [
      { id: 'id', label: 'ID', minWidth: 80, width: 80, align: 'center' },
      { id: 'name', label: 'Name', minWidth: 100,
        format: (val, row) => this.nameCell(val, row) },
      { id: 'created_at', label: 'Created At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
      { id: 'updated_at', label: 'Updated At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
      { id: 'server', label: '', minWidth: 100, width: 100,
        format: (val, row) => this.createServerButton(val, row) },
      { id: 'delete', label: '', minWidth: 100, width: 100,
        format: (val, row) => this.createDeleteButton(val, row) },
    ];

    this.table = React.createRef();
    this.registerDialog = React.createRef();

    this.handleOpenRegister = this.handleOpenRegister.bind(this);
  }

  componentDidMount() {
    //プロジェクトの一覧を取得
    var paging = this.state.paging;
    this.view(paging);

    //同居 server の存在確認(単体起動の ikasbox では 404 になり非表示のまま)
    API.post("/api/v1/server/status").then((res) => {
      this.setState({
        serverMode: true,
        serverProjectId: res.data.projectId
      });
    }).catch(() => {});
  }

  view(paging) {
    var args = {
      paging : paging
    }

    API.post("/api/v1/projects/view", args).then((res) => {
      var result = res.data;
      this.table.current.set(result.projects,result.paging);
      this.setState({
        paging: result.paging
      });
    }).catch((err) => {
      console.log(err)
    });
  }

  contentLink(val, row) {
    return <NameLink href={"/projects/group/" + row["id"]}>{val}</NameLink>;
  }

  createContentsButton(row) {
    return (
      <Button variant="contained" component={RouterLink} to={"/projects/contents/" + row["id"]}>
        Contents
      </Button>
    );
  }

  nameCell(val, row) {
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
        {this.contentLink(val, row)}
        {this.createContentsButton(row)}
      </Box>
    );
  }

  handleDelete(val) {

    Confirm("Project Delete", "Delete project \"" + val.name + "\"?").then(() => {

      var args = {
        projectId: val.id
      }

      API.delete("/api/v1/projects/delete", args).then(() => {
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

  //同居 server の work file をこのプロジェクトで作成する。
  //現在 server に読み込まれているプロジェクトは緑(success)で示す
  handleServerCreate(row) {

    Confirm("Server Create", "Create the server work file for \"" + row.name + "\"?\n" +
      "The running server will switch to this project's contents. " +
      "Run \"ika-client create " + row.id + "\" afterwards to keep the client in sync.").then(() => {

      var args = {
        projectId: row.id
      }

      API.post("/api/v1/server/create", args).then(() => {
        this.setState({ serverProjectId: row.id });
        Alert("Server Create", "Created and reloaded (project " + row.id + ").");
      }).catch((err) => {
        Alert("Server Create", err.message);
      });

    }).catch(() => {});
  }

  createServerButton(val, row) {
    if (!this.state.serverMode) {
      return null;
    }
    var active = row.id === this.state.serverProjectId;
    return (
      <Button color={active ? "success" : "primary"} variant="contained"
        onClick={() => this.handleServerCreate(row)}>
        Server
      </Button>
    );
  }

  handleOpenRegister() {
    this.registerDialog.current.open();
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

  render() {

    var paging = this.state.paging;
    return (<>
      <ProjectRegisterDialog ref={this.registerDialog} onCommit={() => this.view(this.state.paging)} />
      <Button variant="contained" onClick={this.handleOpenRegister}>Register</Button>
      <FlexTable columns={this.columns} paging={paging} ref={this.table}
        onPageChange={this.handleChangePage} onRowsPerPageChange={this.handleChangeRowsPerPage}/>
    </>);
  }
}

export default Projects;