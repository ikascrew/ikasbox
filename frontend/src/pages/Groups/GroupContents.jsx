import React from "react"

import API from "../../API.js"
import Paging from "../Paging.js";

import {withParams} from "../Pages.jsx"
import {TextField,Box,Pagination,Card,CardMedia,CardContent,Typography} from '@mui/material';

class GroupContents extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      group : {},
      contents : [],
      paging : Paging.create(100)
    }

    this.groupId = props.params.group_id;

    this.nameTxt = React.createRef();
    this.pathTxt = React.createRef();
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
      groupId : Number(this.groupId),
      paging : paging
    }

    API.post("/api/v1/groups/contents",args).then( (res) => {
      var result = res.data;

      this.setState({
        group:result.group,
        contents:result.contents,
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

    var group = this.state.group;
    var data = this.state.contents;
    var paging = this.state.paging;

    var pageCount = Math.max(1, Math.ceil(paging.count / paging.limit));

    return (<>

      <TextField
        value={group.name}
        autoFocus margin="dense"
        id="name" type="text" label="Name"
        variant="standard"
        inputRef={this.nameTxt}
        fullWidth
      />

      <TextField
        value={group.path}
        margin="dense"
        id="path" type="text" label="Path"
        variant="standard"
        inputRef={this.pathTxt}
        fullWidth
      />

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
      <Pagination count={pageCount} page={paging.current} onChange={this.handleChangePage} showFirstButton showLastButton />
    </>);
  }

}

export default withParams(GroupContents);
