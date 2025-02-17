import UserBadge from "../../UserBadge/UserBadge.jsx";
import {List, Paper} from "@mui/material";
import {YourUser} from "../../App.jsx"
import ServerBadge from "./ServerBadge/ServerBadge.jsx";
import SearchIcon from '@mui/icons-material/Search';
import "./ServerList.css";


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
        {}
    ]
    
    function searchServer(serverData)
    {
        serverData.preventDefault();
        const data=new FormData(serverData.target);
        const server_ID=formData.get("serverID");
        console.log("Server ID:", server_ID);
    }

    function search(Data)
    {
        Data.preventDefault();
        const query = formData.get("query");
    }

    return(
        <div>
            <Paper elevation={3} sx={{
                borderRadius: 7.5
            }} style={{ display: "flex", flexDirection: "column", width: "15rem", height: "calc(100vh - 2rem)", margin: "1rem", overflow: "hidden"}}>
                <UserBadge user={YourUser.name} status={YourUser.status} online ={YourUser.online} img={YourUser.icon}/>

                <div className="searchData">
                    <form onSubmit={search}>
                        <input name="query" placeholder="Search"/>
                        <button type="submit">
                            <SearchIcon />
                        </button>
                    </form>
                </div>

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