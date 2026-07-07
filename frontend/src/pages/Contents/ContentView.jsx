import React from "react"

import API from "../../API.js"
import {withParams} from "../Pages.jsx"
import {Table,TableBody,TableRow,TableCell} from '@mui/material';

class ContentView extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      content : {}
    }

    this.id = props.params.id;
  }

  componentDidMount() {
    this.view();
  }

  view() {

    var args = {
      id : Number(this.id)
    }

    API.post("/api/v1/contents/view", args).then((res) => {
      var result = res.data;
      this.setState({
        content : result.content
      });
    }).catch((err) => {
      console.log(err)
    });
  }

  render() {

    var content = this.state.content;

    var width = content.width;
    var height = content.height;
    while ( width > 720 ) {
      width = width / 2;
      height = height / 2;
    }

    return (<>

      <video
        src={"/content/media/" + content.id}
        poster={"/thumb/" + content.id}
        width={width}
        height={height}
        controls
        muted
      />

      <Table>
        <TableBody>
          <TableRow>
            <TableCell>ID</TableCell>
            <TableCell>{content.id}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell>{content.name}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>Path</TableCell>
            <TableCell>{content.path}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>Size</TableCell>
            <TableCell>{content.width} x {content.height}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>FPS / Frames</TableCell>
            <TableCell>{content.fps} {content.frames}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell>Fourcc</TableCell>
            <TableCell>{content.fourcc}</TableCell>
          </TableRow>
        </TableBody>
      </Table>

    </>);
  }
}

export default withParams(ContentView);
