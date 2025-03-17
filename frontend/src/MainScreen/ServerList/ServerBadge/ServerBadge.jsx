import {Avatar, Button, Paper, Stack} from "@mui/material";
import {useState} from "react";

function GetUserInfo(){
    alert("I am an alert box!");
}

function ServerBadge(props) {

    const [hovered, setHovered] = useState(false);

    return (
        <Button onClick={GetUserInfo} sx={{
            display: "flex", width: "calc(100% - 1rem)", margin: ".5rem", padding: "0", borderRadius: 10,
            color: "orange", textTransform: "none"}}
                onMouseEnter={() => setHovered(true)}
                onMouseLeave={() => setHovered(false)}
        >
            <Paper elevation={hovered ? 24 : 7} className="UserBadgeContainer" sx={{
                display: "flex", width: "100%", height: "fit-content", margin: "0", borderRadius: 10
            }}>
                <Stack sx={{display: "flex", width: "100%", height: "fit-content", margin: ".5rem"}} direction = "row-reverse" spacing={2}>

                    <Avatar sx={{width: "3rem", height: "3rem"}} src ={props.server?.icon || "serverDefault.png"} alt = "ServerBadgeIcon"></Avatar>

                    <div style={{display: "flex"}} >
                        {props.server?.name || "Server Name"}
                    </div>
                </Stack>
            </Paper>

        </Button>
    )
}

export default ServerBadge;