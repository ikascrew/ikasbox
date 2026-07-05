import React from "react"

import API from "../../API.js"
import Paging from "../Paging.js";

import {Grid,Pagination,Card,CardMedia,CardContent,Typography} from '@mui/material';

class Contents extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      data : [],
      paging : Paging.create()
    }
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
      groupId : -1,
      paging : paging
    }

    API.post("/api/v1/groups/contents",args).then( (res) => {
      var result = res.data;
      this.setState({
        data:result.contents,
        paging:result.paging
      });
    }).catch( (err) => {
      console.log(err)
    });
  }

  handleChangePage = (event, newPage) => {
    var paging = this.state.paging;
    paging.current = newPage;
    this.view(paging);
  };

  render() {

    var data = this.state.data;
    var paging = this.state.paging;

    var pageCount = Math.max(1, Math.ceil(paging.count / paging.limit));

    return (<>

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
      <Pagination count={pageCount} page={paging.current} onChange={this.handleChangePage} showFirstButton showLastButton />
    </>);
  }

}

export default Contents;
