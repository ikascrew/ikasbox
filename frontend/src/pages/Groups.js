import React from "react"

import API from "../API.js"
import Paging from "../Paging.js";

class Groups extends React.Component {

  componentDidMount() {
    //グループの一覧を取得
    var paging = Paging.create(); 
    paging.current = 1;
    this.view(paging);
  }

  view(paging) {
    var args = {
      paging : paging
    }

    API.post("/api/v1/groups/view",args).then( (res) => {
      console.log(res)

    }).catch( (err) => {
      console.log(err)
    });
  }

  render() {
    return <>Groups</>
  }
}

export default Groups;