import React, { useState, useEffect } from 'react';
import { Paper, TextField, Button, Typography, Box } from '@mui/material';
import {useUser} from "../UserContext.jsx";

async function FetchUserData(username, setUserData) {
    try {
        const response = await fetch("http://localhost:3000/v1/api/account/search", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username }),
        });

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const jsonData = await response.json();
        console.log("Full JSON:", jsonData);
        setUserData(jsonData[0]); // ✅ Update state with response data
    } catch (error) {
        console.error("Error:", error.message);
    }
}

function Login({ onLogin }) {
    const { setUser } = useUser();
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [userData, setLocalUserData] = useState(null); // ✅ Local state to track user data

    const handleLogin = async () => {
        await FetchUserData(username, setLocalUserData);

        setUser(userData);

    };

    useEffect(() => {
        console.log("userData", userData);

        if (userData && userData.id != null) {
            setUser(userData);  // Set user data using `setUser` from the context
            onLogin(true);  // Trigger login state change in the parent component
        }
    }, [userData, onLogin]);

    return (
        <Box
            sx={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                height: '100vh',
                background: 'linear-gradient(45deg, hsla(29, 100%, 57%, 1) 0%, hsla(0, 100%, 79%, 1) 100%)'
            }}
        >
            <Paper
                elevation={10}
                sx={{
                    padding: 4,
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 2,
                    width: '300px',
                    borderRadius: 3
                }}
            >
                <Typography variant="h4" align="center">
                    Login
                </Typography>
                <TextField
                    label="Username"
                    variant="outlined"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    fullWidth
                />
                <TextField
                    label="Password"
                    type="password"
                    variant="outlined"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    fullWidth
                />
                <Button
                    variant="contained"
                    color="primary"
                    onClick={handleLogin}
                    fullWidth
                >
                    Sign In
                </Button>
                <Typography variant="body2" align="center">
                    Don't have an account? Register
                </Typography>
            </Paper>
        </Box>
    );
}

export default Login;
