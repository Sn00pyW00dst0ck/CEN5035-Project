import UserBadge from "../../UserBadge/UserBadge.jsx";
import {List, Paper, TextField} from "@mui/material";
import {YourUser} from "../../App.jsx"
import ServerBadge from "./ServerBadge/ServerBadge.jsx";
import "./ServerList.css";
import {useState} from "react";

function searchServer(serverData)
    {
        serverData.preventDefault();
        const data=new FormData(serverData.target);
        const server_ID=formData.get("serverID");
        console.log("Server ID:", server_ID);
    }




function ServerList() {

    const [query, setQuery] = useState("");

    function handleServerSearch(event) {
        event.preventDefault();
        setQuery(event.target.value);
    }

    const servers = [
        { id: 1, name: "test1", icon: "public/vite.svg"},
        { id: 2, name: "test2", icon: "public/vite.svg"},
        { id: 3, name: "Test1", icon: "public/vite.svg"},
        { id: 4, name: "Test2", icon: "public/vite.svg"},
        { id: 5, name: "thisIsATest1", icon: "public/vite.svg"},
        { id: 6, name: "Alice", icon: "public/vite.svg"}
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

                    <TextField id="serverSearchInput" sx ={{display:"flex", margin: ".5rem"}}
                               label="Search Server"
                               type="search"
                               size="small"
                               onChange={handleServerSearch}
                    />

                    <div id='serverBadgeHolder'>
                        {servers
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