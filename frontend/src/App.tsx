import React from 'react';
import './App.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import HomePage from './pages/HomePage';
import Layout from './components/Layout';
import { createTheme, ThemeProvider } from '@mui/material/styles';

let theme = createTheme({
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          color: "#641803",
          fontWeight: 700,
        }
      }
    }
  },
  typography: {
    fontFamily: 'Times New Roman, serif', // Set the font family to Times New Roman
    fontSize: 20, // Set the base font size
    fontWeightRegular: 1000, // Set the regular font weight
    fontWeightBold: 1000, // Set the bold font weight
  },
});

function App() {
  return (
    <ThemeProvider theme={theme}>
      <Layout>
        <BrowserRouter>
          <Routes>
            <Route path="/">
              <Route index element={<HomePage />} />
            </Route>
          </Routes>
        </BrowserRouter>
      </Layout>
    </ThemeProvider>
  );
}

export default App;
