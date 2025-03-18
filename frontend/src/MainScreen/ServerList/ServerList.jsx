import UserBadge from "../../UserBadge/UserBadge.jsx";
import { List, Paper } from "@mui/material";
import { YourUser } from "../../App.jsx";
import ServerBadge from "./ServerBadge/ServerBadge.jsx";
import "./ServerList.css";
import { useState } from "react";
import Search from "../../CommonComponents/Search/Search.jsx";

function searchServer(event) {
    event.preventDefault();
    const data = new FormData(event.target);
    const server_ID = data.get("serverID");
    console.log("Server ID:", server_ID);
}

function ServerList({ servers }) {
    const [query, setQuery] = useState("");
    const [selectedServer, setSelectedServer] = useState(null);

    function handleServerSearch(event) {
        event.preventDefault();
        setQuery(event.target.value);
    }

    function handleServerClick(server) {
        console.log("Selected Server:", server);
        setSelectedServer(server); // ✅ Updates with correct server & channels
    }

    return (
        <div style={{ display: "flex" }}>
            {/* Server List */}
            <Paper elevation={3} sx={{
                borderRadius: 7.5,
                display: "flex",
                flexDirection: "column",
                width: "15rem",
                height: "calc(100vh - 2rem)",
                margin: "1rem",
                overflow: "hidden"
            }}>
                <UserBadge user={YourUser.name} status={YourUser.status} online={YourUser.online} img={YourUser.icon} />

                <List sx={{ display: "flex", flexDirection: "column", width: '100%', height: "100%", overflow: 'auto' }}>
                    <Search id="serverSearchInput" return={handleServerSearch} />

                    <div id='serverBadgeHolder'>
                        {servers
                            .filter((server) => server.name.toLowerCase().includes(query.toLowerCase()))
                            .map((server) => (
                                <li key={server.id} onClick={() => handleServerClick(server)} style={{ cursor: "pointer" }}>
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

            {/* Channel List - Only Shows When a Server is Selected */}
            {selectedServer && selectedServer.channels && (
                <Paper elevation={3} sx={{
                    borderRadius: 7.5,
                    display: "flex",
                    flexDirection: "column",
                    width: "15rem",
                    height: "calc(100vh - 2rem)",
                    margin: "1rem",
                    overflow: "hidden"
                }}>
                    <h3 style={{ textAlign: "center", marginTop: "1rem" }}>
                        {selectedServer.name} Channels
                    </h3>
                    <List sx={{ overflow: 'auto' }}>
                        {selectedServer.channels.map((channel, index) => (
                            <li key={`${selectedServer.id}-channel-${index}`}>
                                <ServerBadge server={{ id: `${selectedServer.id}-channel-${index}`, name: channel }} />
                            </li>
                        ))}
                    </List>
                </Paper>
            )}
        </div>
    );
}

export default ServerList;