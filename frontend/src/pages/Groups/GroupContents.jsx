import React from "react"

import API from "../../API.js"
import Paging from "../Paging.js";

import {withParams} from "../Pages.jsx"
import {TextField,Box,Stack,Pagination,Card,CardMedia,CardContent,Typography,Button} from '@mui/material';
import LoadingButton from "../../components/LoadingButton.jsx";
import ContentRegisterDialog from "./ContentRegisterDialog.jsx";
import { Alert } from "../Layout.jsx";
import { Link as RouterLink } from "react-router";

class GroupContents extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      group : {},
      contents : [],
      paging : Paging.create(100),
      name : ""
    }

    this.groupId = props.params.group_id;
    this.registerDialog = React.createRef();
  }

  handleOpenRegister = () => {
    this.registerDialog.current.open();
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
        paging:result.paging,
        name:result.group.name
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

  handleChangeName = (event) => {
    this.setState({
      name: event.target.value
    });
  };

  handleSaveName = () => {

    var args = {
      groupId: Number(this.groupId),
      name: this.state.name
    }

    return API.patch("/api/v1/groups/rename", args).then(() => {
      this.view(this.state.paging);
      Alert("Group", "Name updated.");
    }).catch((err) => {
      console.log(err);
    });
  };

  render() {

    var group = this.state.group;
    var data = this.state.contents;
    var paging = this.state.paging;

    var pageCount = Math.max(1, Math.ceil(paging.count / paging.limit));

    return (<>

      <Stack direction="row" spacing={2} alignItems="center">
        <TextField
          value={this.state.name}
          onChange={this.handleChangeName}
          autoFocus margin="dense"
          id="name" type="text" label="Name"
          variant="standard"
          InputLabelProps={{ shrink: true }}
          fullWidth
        />
        <LoadingButton onClick={this.handleSaveName}>Save</LoadingButton>
      </Stack>

      <TextField
        value={group.path || ""}
        margin="dense"
        id="path" type="text" label="Path"
        variant="standard"
        InputLabelProps={{ shrink: true }}
        InputProps={{ readOnly: true }}
        fullWidth
      />

      <Stack direction="row" justifyContent="flex-end" sx={{ marginTop: 1 }}>
        <ContentRegisterDialog
          ref={this.registerDialog}
          groupId={this.groupId}
          onCommit={() => {
            this.view(this.state.paging);
            Alert("Content", "Registered.");
          }}
        />
        <Button variant="contained" onClick={this.handleOpenRegister}>Add Content</Button>
      </Stack>

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

export default withParams(GroupContents);
