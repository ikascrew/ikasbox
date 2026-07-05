import React from "react"

import API from "../../API";
import {withParams} from "../Pages.jsx"
import {Grid,Card,CardMedia,CardContent,Typography} from '@mui/material';

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

      <Grid container sx={{marginTop:"10px"}} justify="center">

{data.map( (obj) => {

  return (
    <Grid item container xs={12} md={4} sm={6} spacing={1} key={"content-" + obj.id}>
      <Card sx={{ maxWidth:300, minWidth: 300 ,maxHeight:300,margin:5}}>
        <a href={"/contents/" + obj.id}>
        <CardMedia
          sx={{ height: 200 }}
          image={"/thumb/" + obj.id}
        />
        </a>

        <CardContent>
          <Typography gutterBottom variant="h6" component="div"> {obj.name} </Typography>
          <Typography variant="body2" color="text.secondary"></Typography>
        </CardContent>

      </Card>
    </Grid>
  );
})}

      </Grid>
    </>);
  }
}

export default withParams(ProjectContents);
