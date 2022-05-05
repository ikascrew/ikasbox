import React from "react";
import { createRoot } from 'react-dom/client';

import Layout from "./Layout.js"
import "./css/App.css"

class App extends React.Component {
  render() {
    return (<>
      <Layout/>
    </>);
  }
}

const container = document.getElementById('app');
const root = createRoot(container);
root.render(<App />);

