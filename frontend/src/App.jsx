import React from "react";
import { createRoot } from 'react-dom/client';
import { BrowserRouter } from "react-router";
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';

import theme from "./theme.js";
import Layout from "./pages/Layout.jsx"

class App extends React.Component {
  render() {
    return (<>
      <ThemeProvider theme={theme}>
        <CssBaseline />
        <BrowserRouter>
          <Layout/>
        </BrowserRouter>
      </ThemeProvider>
    </>);
  }
}

const container = document.getElementById('app');
const root = createRoot(container);
root.render(<App />);

