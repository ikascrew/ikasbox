import React from "react"

import API from "../../API";
import {withParams} from "../Pages.jsx"
import {Box,Card,CardMedia,CardContent,Typography} from '@mui/material';

class ProjectContents extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      project : {},
      contents : []
    }

    this.projectId = this.props.params.id;
  }

  componentDidMount() {
    this.view();
  }

  view() {

    var args = {
      projectId : Number(this.projectId)
    }

    API.post("/api/v1/projects/contents", args).then((res) => {
      var result = res.data;
      this.setState({
        project : result.project,
        contents : result.contents
      });
    }).catch((err) => {
      console.log(err)
    });
  }

  render() {

    var project = this.state.project;
    var data = this.state.contents;

    return (<>

      <Typography variant="h5" sx={{marginTop:"10px"}}>{project.name}</Typography>

      <Box sx={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fill, minmax(140px, 1fr))",
        gap: 1.5,
        marginTop: "10px",
      }}>

{data.map( (obj) => {

  return (
    <Card key={"content-" + obj.id}>
      <a href={"/contents/" + obj.id}>
        <CardMedia
          sx={{ height: 100 }}
          image={"/thumb/" + obj.id}
        />
      </a>

      <CardContent sx={{ padding: 1, "&:last-child": { paddingBottom: 1 } }}>
        <Typography variant="body2" component="div" noWrap> {obj.name} </Typography>
      </CardContent>

    </Card>
  );
})}

      </Box>
    </>);
  }
}

export default withParams(ProjectContents);
