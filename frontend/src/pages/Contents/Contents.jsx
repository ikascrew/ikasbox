import React from "react"

import API from "../../API.js"
import Paging from "../Paging.js";

import {Box,Pagination,Card,CardMedia,CardContent,Typography} from '@mui/material';
import { Link as RouterLink } from "react-router";

class Contents extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      data : [],
      paging : Paging.create(100)
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

      <Box sx={{
        display: "grid",
        gridTemplateColumns: "repeat(auto-fill, minmax(140px, 1fr))",
        gap: 1.5,
        marginTop: "10px",
      }}>

{data.map( (obj) => {

  return (
    <Card key={"content-" + obj.id}>
      <RouterLink to={"/contents/" + obj.id}>
        <CardMedia
          sx={{ height: 100 }}
          image={"/thumb/" + obj.id}
        />
      </RouterLink>

      <CardContent sx={{ padding: 1, "&:last-child": { paddingBottom: 1 } }}>
        <Typography variant="body2" component="div" noWrap> {obj.name} </Typography>
      </CardContent>

    </Card>
  );
})}

      </Box>
      <Pagination count={pageCount} page={paging.current} onChange={this.handleChangePage} showFirstButton showLastButton />
    </>);
  }

}

export default Contents;
