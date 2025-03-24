import './App.css'
import { ThemeProvider, createTheme } from '@mui/material/styles';
import MainScreen from "./MainScreen/MainScreen.jsx";
import Login from "./Login/Login.jsx";
import { useState } from 'react';

export const YourUser = {id: 1, name: "YourUsername", status: "Hi", online: true, icon: "public/vite.svg"}

const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
const theme = createTheme({
  palette: {
    mode: darkThemeMq.matches ? 'dark' : 'light',
  },
});

function App() {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  return (
    <div>
      <ThemeProvider theme={theme}>
        {isLoggedIn ? <MainScreen /> : <Login onLogin={setIsLoggedIn} />}
      </ThemeProvider>
    </div>
  )
}

export default App