import React from "react";

class Paging {
  static create() {
    return {
      current:1,
      count:0,
      limit:10,
      maxPage:0
    }
  }
}

export default Paging;