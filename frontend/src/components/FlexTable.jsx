import React from "react"

import {
  Paper, TableContainer, Table, TableHead, TableRow, TableCell, TableBody,
  TablePagination
} from '@mui/material';
import Paging from "../pages/Paging";

class FlexTable extends React.Component {

  constructor(props) {
    super(props);

    var columns = props.columns;
    var paging = props.paging;
    //var data = props.data;

    this.state = {
      columns : columns,
      data : [],
      paging : paging
    }
  }

  set(data,paging) {
    this.setState({
      data:data,
      paging: paging
    });
  }

  componentDidUpdate() {
  }

  handleChangePage = (event, newPage) => {
    if ( this.props.onPageChange !== undefined ) {
      this.props.onPageChange(event, newPage + 1);
    }
  }

  handleChangeRowsPerPage = (event) => {
    if ( this.props.onRowsPerPageChange !== undefined ) {
      this.props.onRowsPerPageChange(event);
    }
  }

  render () {

    let columns = this.state.columns;
    let data = this.state.data;
    let paging = this.state.paging;

    return (<>
      <Paper sx={{ width: '100%', overflow: 'hidden', marginTop: "10px" }}>
        <TableContainer sx={{ maxHeight: 1000 }}>
          <Table stickyHeader aria-label="sticky table">
            <TableHead>
              <TableRow>
                {columns.map((column) => (
                  <TableCell
                    key={column.id} align={column.align}
                    style={{ minWidth: column.minWidth, width: column.width }}
                  >
                    {column.label}
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>

            <TableBody>
              {data.map((row, idx) => {
                return (
                  <TableRow hover tabIndex={-1} key={row.id + "-" + idx}>
                    {columns.map((column, idx) => {
                      const value = row[column.id];
                      return (
                        <TableCell key={column.id + "-" + idx} align={column.align}
                          style={{ minWidth: column.minWidth, width: column.width }}
                        >
                          {column.format ? column.format(value, row) : value}
                        </TableCell>
                      );
                    })}
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
{ paging !== undefined && 
        <TablePagination
          rowsPerPageOptions={[10, 25, 100]}
          component="div"
          count={paging.count}
          rowsPerPage={paging.limit}
          page={paging.current - 1}
          onPageChange={this.handleChangePage}
          onRowsPerPageChange={this.handleChangeRowsPerPage}
        />
}
      </Paper>
    </>);
  }
}

export default FlexTable