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
import { Confirm } from "../Layout.jsx";

class Projects extends React.Component {
  
  constructor(props) {
    super(props);

    this.state = {
      data : [],
      paging : Paging.create()
    }

    this.columns = [
      { id: 'id', label: 'ID', minWidth: 80, width: 80, align: 'center' },
      { id: 'name', label: 'Name', minWidth: 100,
        format: (val, row) => this.nameCell(val, row) },
      { id: 'created_at', label: 'Created At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
      { id: 'updated_at', label: 'Updated At', minWidth: 190, width: 190, align: 'center',
        format: (value) => Util.formatDate(value) },
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

  handleOpenContents(row) {
    location.href = "/projects/contents/" + row["id"];
  }

  createContentsButton(row) {
    return (
      <Button variant="contained" onClick={() => this.handleOpenContents(row)}>
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

      API.delete("/api/v1/projects/delete", { data: args }).then(() => {
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