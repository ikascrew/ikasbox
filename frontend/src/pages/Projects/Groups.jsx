import React from "react"

import Select from "../../components/FormSelect";
import {withParams} from "../Pages.jsx"

import API from "../../API";
import Util from "../../Util";

import {
  Button
} from '@mui/material';

import FlexTable from "../../components/FlexTable";

class ProjectGroups extends React.Component {
  
  constructor(props) {
    super(props);
    this.state = {
      groups : [],
      value: ""
    }

    this.columns = [
      { id: 'id', label: 'ID', minWidth: 20 },
      { id: 'name', label: 'Name', minWidth: 100,
        format: (val, row) => this.contentLink(val, row) },
      { id: 'created_at', label: 'Created At', minWidth: 80,
        format: (value) => Util.formatDate(value) },
      { id: 'updated_at', label: 'Updated At', minWidth: 80,
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

      if ( this.groups === undefined ) {
        let groups = new Map();
        result.allGroups.forEach( (elm) => {
          groups.set(elm.id,elm.name);
        })

        this.setState({
          groups : groups
        });
      }

      this.table.current.set(result.groups);

    }).catch((err) => {
      console.log(err)
    });
  }

  contentLink(val, row) {
    return <a href={"/groups/contents/" + row["id"]}>{val}</a>;
  }

  handleChangeValue = (groupId) => {

    var args = {
      projectId : Number(this.projectId),
      groupId : Number(groupId)
    }

    API.patch("/api/v1/projects/group/add", args).then((res) => {
      this.setState({ value: "" });
      this.view();
    }).catch((err) => {
      console.log(err)
    });
  }

  render() {

    return (<>
      <Select value={this.state.value} values={this.state.groups} onChange={this.handleChangeValue} empty="Add group..." />
      <FlexTable columns={this.columns} ref={this.table} />
    </>);
  }
}

export default withParams(ProjectGroups);