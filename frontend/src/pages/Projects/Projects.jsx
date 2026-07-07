import React from "react"
import Paging from "../Paging";
import API from "../../API";
import Util from "../../Util";

import {
  Button
} from '@mui/material';

import FlexTable from "../../components/FlexTable";
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
      { id: 'id', label: 'ID', minWidth: 20 },
      { id: 'name', label: 'Name', minWidth: 100,
        format: (val, row) => this.contentLink(val, row) },
      { id: 'created_at', label: 'Created At', minWidth: 80,
        format: (value) => Util.formatDate(value) },
      { id: 'updated_at', label: 'Updated At', minWidth: 80,
        format: (value) => Util.formatDate(value) },
      { id: 'contents', label: '', minWidth: 100,
        format: (val, row) => this.contentsLink(row) },
      { id: 'delete', label: '', minWidth: 100,
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
    return <a href={"/projects/group/" + row["id"]}>{val}</a>;
  }

  contentsLink(row) {
    return <a href={"/projects/contents/" + row["id"]}>Contents</a>;
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