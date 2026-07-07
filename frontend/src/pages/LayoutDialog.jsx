import React from "react";

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';

import LoadingButton from "../components/LoadingButton";

class LayoutDialog extends React.Component {

  constructor(props) {
    super(props);

    let loading = undefined;
    if ( props.loading !== undefined ) {
      loading = props.loading;
    }

    this.state = {
      view:false,
      type:"alert",
      loading:loading,
      title:"",
      message:""
    }

    this.promise = undefined;

    this.handleYes = this.handleYes.bind(this);
    this.handleNo = this.handleNo.bind(this);
    this.handleLodingYes = this.handleLoadingYes.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.close = this.close.bind(this);
  }

  handleClose() {
    this.resolv("Close");
    this.close();
  }

  handleLoadingYes() {
    this.resolv("Loading Yes");
    let p = new Promise( (resolv,reject) => {
      this.state.loading(resolv,reject);
    });
    p.finally( () => {
      this.close();
    });
    return p;
  }

  handleYes() {
    this.resolv("Yes");
    this.close();
  }

  handleNo() {
    this.reject("No");
    this.close();
  }

  close() {
    this.setState({
      view:false
    });
  }

  show(title,msg,type) {

    this.setState({
      view:true,
      title:title,
      message:msg,
      type : type
    });

    var self = this;
    this.promise = new Promise( (resolv,reject) => {
      self.resolv = resolv;
      self.reject = reject;
    });
    return this.promise;
  }

  render() {
    return (<>
      <Dialog open={this.state.view}>
        <DialogTitle>{this.state.title}</DialogTitle>
        <DialogContent>
          <DialogContentText sx={{ whiteSpace: "pre-line" }}>{this.state.message}</DialogContentText>
        </DialogContent>
        <DialogActions>
{this.state.type === "alert" &&
<>
          <Button onClick={this.handleClose}>Close</Button>
</>
}
{this.state.type === "confirm" &&
<>
          <Button onClick={this.handleNo}>No</Button>
          <Button color="success" variant="contained" onClick={this.handleYes}>Yes</Button>
</>
}
        </DialogActions>
      </Dialog>
    </>);
  }
}


export default LayoutDialog;