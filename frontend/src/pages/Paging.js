import React from "react";

class Paging {
  static create(limit = 10) {
    return {
      current:1,
      count:0,
      limit:limit
    }
  }
}

export default Paging;