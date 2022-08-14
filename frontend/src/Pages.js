import React from "react";

import { Routes,Route } from "react-router-dom";
import Index from "./pages/Index.js"
import Groups from "./pages/Groups.js"
import Projects from "./pages/Projects.js"

class Pages extends React.Component {
  render() {
    return (
      <Routes>
        <Route path="/" element={<Index/>} />
        <Route path="/groups" element={<Groups/>}/>
        <Route path="/projects" element={<Projects/>}/>
      </Routes>
    )
  }
}

export default Pages;