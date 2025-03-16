import UserBadge from "../../UserBadge/UserBadge.jsx";
import {List, Paper, TextField} from "@mui/material";
import {YourUser} from "../../App.jsx"
import ServerBadge from "./ServerBadge/ServerBadge.jsx";
import "./ServerList.css";
import {useState} from "react";
import Search from "../../CommonComponents/Search/Search.jsx";

function searchServer(serverData)
    {
        serverData.preventDefault();
        const data=new FormData(serverData.target);
        const server_ID=formData.get("serverID");
        console.log("Server ID:", server_ID);
    }




function ServerList(props) {

    const [query, setQuery] = useState("");

    function handleServerSearch(event) {
        event.preventDefault();
        setQuery(event.target.value);
    }

    return(
        <div>
            <Paper elevation={3} sx={{
                borderRadius: 7.5,
                display: "flex",
                flexDirection: "column",
                width: "15rem",
                height: "calc(100vh - 2rem)",
                margin: "1rem",
                overflow: "hidden"
            }}>
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

                    <Search id = "serverSearchInput" return = {handleServerSearch}/>

                    <div id='serverBadgeHolder'>
                        {props.servers
                            .filter((server) => server.name.toLowerCase().includes(query.toLowerCase()))
                            .map((server) => (
                                <li key={`section-${server.id}`}>
                                    <ServerBadge server={server} />
                                </li>
                            ))}
                    </div>

                </List>

                <div className="joinServer">
                    <form onSubmit={searchServer}>
                        <input name="serverID" placeholder="Enter a Server ID" />
                        <button type="submit">Join</button>
                    </form>
                </div>

            </Paper>
        </div>
    )
}

export default ServerList;