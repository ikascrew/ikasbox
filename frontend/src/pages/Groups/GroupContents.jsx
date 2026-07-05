import React from "react"

import API from "../../API.js"
import Paging from "../Paging.js";

import {withParams} from "../Pages.jsx"
import {TextField,Grid,Pagination,Card,CardMedia,CardContent,Typography,CardActions,Button} from '@mui/material';
import GroupImportDialog from "./GroupImportDialog.jsx";

class GroupContents extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      group : {},
      contents : [],
      paging : Paging.create()
    }

    this.groupId = props.params.group_id;

    this.nameTxt = React.createRef();
    this.pathTxt = React.createRef();
    this.importDialog = React.createRef();
  }

  handleOpenImport = () => {
    this.importDialog.current.open();
  };

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

      <GroupImportDialog ref={this.importDialog} groupId={this.groupId} />
      <Button variant="contained" sx={{marginTop:"10px"}} onClick={this.handleOpenImport}>Import</Button>

      <Grid container sx={{marginTop:"10px"}} justify="center">

{data.map( (obj) => {

  return (
    <Grid item container xs={12} md={4}  sm={6} spacing={1} key={"content-" + obj.id}>
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

export default withParams(GroupContents);
