import React from "react";

import { useParams,Routes,Route } from "react-router-dom";
import Index from "./Index.jsx"
import Groups from "./Groups/Groups.jsx"
import Contents from "./Contents/Contents.jsx"
import ContentView from "./Contents/ContentView.jsx"
import GroupContents from "./Groups/GroupContents.jsx"
import Projects from "./Projects/Projects.jsx"
import ProjectGroups from "./Projects/Groups.jsx"
import ProjectContents from "./Projects/ProjectContents.jsx"

class Pages extends React.Component {
  render() {
    return (
      <Routes>
        <Route path="/" element={<Index/>} />
        <Route path="/groups" element={<Groups/>}/>
        <Route path="/projects" element={<Projects/>}/>
        <Route path="/projects/group/:id" element={<ProjectGroups/>}/>
        <Route path="/projects/contents/:id" element={<ProjectContents/>}/>
        <Route path="/contents" element={<Contents/>}/>
        <Route path="/contents/:id" element={<ContentView/>}/>
        <Route path="/groups/contents/:group_id" element={<GroupContents/>}/>
      </Routes>
    )
  }
}

export function withParams(Component) {
  return props => <Component {...props} params={useParams()} />;
}

export default Pages;