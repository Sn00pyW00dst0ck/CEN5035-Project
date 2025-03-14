import './App.css'

import { ThemeProvider, createTheme } from '@mui/material/styles';
import MainScreen from "./MainScreen/MainScreen.jsx";

//Temporary should be fixed.
export const YourUser = {id: 1, name: "Alice", status: "Hi", online: true, icon: "public/vite.svg"}

const darkThemeMq = window.matchMedia("(prefers-color-scheme: dark)");

const theme = createTheme({
    palette: {
        mode: darkThemeMq.matches ? 'dark' : 'light',
    },
});

function App() {

    return (
        <div>
            <ThemeProvider theme={theme}>
                <MainScreen/>
            </ThemeProvider>
        </div>
    )
}

export default App

