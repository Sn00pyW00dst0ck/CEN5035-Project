import './App.css'

import { ThemeProvider, createTheme } from '@mui/material/styles';
import MainScreen from "./MainScreen/MainScreen.jsx";

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

function App() {

  return (
      <div>
        <ThemeProvider theme={darkTheme}>
          <MainScreen/>
        </ThemeProvider>
      </div>
  )
}

export default App
