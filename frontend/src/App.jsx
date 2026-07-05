import React from "react";
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from "react-router-dom";

import Layout from "./pages/Layout.jsx"
import "./css/App.css"

class App extends React.Component {
  render() {
    return (<>
      <BrowserRouter>
        <Layout/>
      </BrowserRouter>
    </>);
  }
}

const container = document.getElementById('app');
const root = createRoot(container);
root.render(<App />);

