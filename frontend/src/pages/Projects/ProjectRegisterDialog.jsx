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

class ProjectRegisterDialog extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      view : false
    }

    this.handleRegister = this.handleRegister.bind(this);

    this.commitFunc = props.onCommit;
    this.nameTxt = React.createRef();
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

  handleRegister = () => {

    var name = this.nameTxt.current.value;
    return new Promise( (resolv,reject) => {

      let args = {
        name : name
      }

      API.patch("/api/v1/projects/register",args).then( (res) => {
        if ( this.commitFunc !== undefined ) {
          this.commitFunc();
        }
        this.handleClose();
        resolv("success");
      }).catch( (err) => {
        reject("error");
      });
    });

  };

  render () {
    return (<>
      <Dialog open={this.state.view}>
        <DialogTitle>Register</DialogTitle>
        <DialogContent>
          <DialogContentText>
          </DialogContentText>

          <TextField
            autoFocus margin="dense"
            id="name" type="text" label="Name"
            variant="standard"
            inputRef={this.nameTxt}
            fullWidth
          />

        </DialogContent>

        <DialogActions>
          <Button onClick={this.handleClose}>Cancel</Button>
          <LoadingButton  onClick={this.handleRegister}>Register</LoadingButton>
        </DialogActions>
      </Dialog>
    </>);
  }
}

export default ProjectRegisterDialog;