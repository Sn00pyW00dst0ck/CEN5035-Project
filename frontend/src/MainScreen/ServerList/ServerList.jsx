import UserBadge from "../../UserBadge/UserBadge.jsx";
import {List, Paper} from "@mui/material";
import {YourUser} from "../../App.jsx"
import ServerBadge from "./ServerBadge/ServerBadge.jsx";

function ServerList() {

    const servers = [
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {}
    ]

    return(
        <div>
            <Paper elevation={3} sx={{
                borderRadius: 7.5
            }} style={{ display: "flex", flexDirection: "column", width: "15rem", height: "calc(100vh - 2rem)", margin: "1rem", overflow: "hidden"}}>
                <UserBadge user={YourUser.name} status={YourUser.status} online ={YourUser.online} img={YourUser.icon}/>

                <List
                    sx={{
                        display: "flex",
                        flexDirection: "column",
                        width: '100%',
                        height: "100%",
                        position: 'relative',
                        overflow: 'auto',

                        '& ul': { padding: 0 },
                    }}
                    subheader={<li />}
                >
                    {servers.map((server) => (
                        <li key={`section-${server}`}>
                            <ServerBadge />
                        </li>
                    ))}
                </List>

            </Paper>
        </div>
    )
}

export default ServerList;