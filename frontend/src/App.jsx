import './App.css';
import { UserProvider } from './UserContext';
import { ThemeProvider, createTheme } from '@mui/material/styles';
import MainScreen from "./MainScreen/MainScreen.jsx";
import Login from "./Login/Login.jsx";
import {useState} from 'react';

// Detect system theme preference
const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");
const theme = createTheme({
    palette: {
        mode: darkThemeMq.matches ? 'dark' : 'light',
    },
});

function App() {
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    return (
        <UserProvider>
            <ThemeProvider theme={theme}>
                {isLoggedIn ? (
                    <MainScreen /> // ✅ Pass userData to MainScreen
                ) : (
                    <Login onLogin={setIsLoggedIn} /> // ✅ Fix prop name
                )}
            </ThemeProvider>
        </UserProvider>
    );
}

export default App;
