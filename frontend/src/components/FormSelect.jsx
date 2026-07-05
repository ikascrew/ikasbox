import React from "react";

import {
  FormControl,InputLabel,Select,MenuItem
} from '@mui/material';

class FromSelect extends React.Component {

  constructor(props) {
    super(props);
    this.state = {
      value : "",
      values : new Map()
    }
  }

  handleChangeValue = (e) => {
    let onChange = this.props.onChange;
    if ( onChange !== undefined ) {
      onChange(e.target.value);
    }
  }

  render() {

    let props = this.props;
    let values = props.values;

    if ( values === undefined ) {
      values = new Map();
    } else if ( !(values instanceof Map) ) {
      console.warn("values instance is not Map.");
      values = new Map();
    }

    return (<>
      <FormControl sx={props.sx} fullWidth>

        <InputLabel id={props.id}>{props.label}</InputLabel>

        <Select
          labelId={props.id}
          label={props.label}
          value={props.value}
          onChange={this.handleChangeValue}
        >

        {props.empty !== undefined &&
          <em>
            <MenuItem value="">{props.empty}</MenuItem>
          </em>
        }

        {Array.from(values.entries()).map( ([key,val]) => {
          return (
            <MenuItem key={key} value={key}>{val}</MenuItem>
          );
        })}

        </Select>
      </FormControl>
    </>);
  }
}

export default FromSelect;