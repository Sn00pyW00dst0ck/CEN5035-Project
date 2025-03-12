import * as React from 'react';
import MenuIcon from '@mui/icons-material/Menu';
import {Button, Paper, TextField} from "@mui/material";
import "./ActiveServer.css";
import Search from "../../../CommonComponents/Search/Search.jsx";
import ServerBadge from "../../ServerList/ServerBadge/ServerBadge.jsx";
import MenuBar from "./MenuBar/MenuBar.jsx";

function ActiveServer({setVisible}) {

    return (
        <div style={{display: "flex", flexDirection: "column", height: "100%", width: "100%"}}>

            <MenuBar setVisible={setVisible}/>

            <Paper
                elevation={7}
                sx={{
                    borderRadius: 7.5,
                    borderTopLeftRadius: 0,
                    borderTopRightRadius: 0,
                    borderBottomRightRadius: 0,
                }}
                className="text"
            >
                <TextField
                    sx={{
                        width: "100%",
                        margin: ".5rem",
                        marginLeft: "1rem",
                    }}
                    placeholder="Text Message"
                />
                <button className="sendButton">Send</button>
            </Paper>
        </div>
    );
}

export default ActiveServer;