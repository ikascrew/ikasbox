import React from "react";

import TextField from '@mui/material/TextField';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

import LoadingButton from "../../components/LoadingButton";
import API from "../../API";
import { Alert } from "../Layout.jsx";

class GroupImportDialog extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      view : false
    }

    this.groupId = props.groupId;
    this.pathTxt = React.createRef();
  }

  open = () => {
    this.setState({
      view:true
    });
  };

  handleClose = () => {
    this.setState({
      view:false
    });
  };

  handleImport = () => {

    var path = this.pathTxt.current.value;

    return new Promise( (resolv,reject) => {

      let args = {
        groupId : Number(this.groupId),
        path : path
      }

      API.patch("/api/v1/groups/import",args).then( (res) => {
        this.handleClose();
        Alert("Import", "Import started in the background. Contents will appear as they finish processing.");
        resolv("success");
      }).catch( (err) => {
        reject("error");
      });
    });

  };

  render () {
    return (<>
      <Dialog open={this.state.view}>
        <DialogTitle>Import</DialogTitle>
        <DialogContent>
          <DialogContentText>
          Scan a directory on the server for media files and register them to this group.
          Processing (thumbnail generation) runs in the background; this dialog will not wait for it to finish.
          </DialogContentText>

          <TextField
            autoFocus margin="dense"
            id="path" type="text" label="Directory Path"
            variant="standard"
            inputRef={this.pathTxt}
            fullWidth
          />

        </DialogContent>

        <DialogActions>
          <Button onClick={this.handleClose}>Cancel</Button>
          <LoadingButton onClick={this.handleImport}>Import</LoadingButton>
        </DialogActions>
      </Dialog>
    </>);
  }
}

export default GroupImportDialog;
