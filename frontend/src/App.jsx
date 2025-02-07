import './App.css'

import { ThemeProvider, createTheme } from '@mui/material/styles';
import MainScreen from "./MainScreen/MainScreen.jsx";

//Temporary should be fixed.
export const YourUser = {id: 1, name: "Alice", status: "Hi", online: true, icon: "public/vite.svg"}

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
