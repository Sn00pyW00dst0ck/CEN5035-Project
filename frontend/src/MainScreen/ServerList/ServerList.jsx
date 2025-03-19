import UserBadge from "../../UserBadge/UserBadge.jsx";
import { List, Paper, TextField } from "@mui/material";
import { YourUser } from "../../App.jsx";
import ServerBadge from "./ServerBadge/ServerBadge.jsx";
import Search from "../../CommonComponents/Search/Search.jsx";
import { useState } from "react";
import "./ServerList.css";

function ServerList({ servers }) {
  const [query, setQuery] = useState("");
  const [selectedServer, setSelectedServer] = useState(null);
  const [newChannelName, setNewChannelName] = useState("");
  const [showAddChannelForm, setShowAddChannelForm] = useState(false);

  function handleServerSearch(event) {
    setQuery(event.target.value);
  }

  function handleServerClick(server) {
    console.log("Selected Server:", server);
    setSelectedServer(server);
  }

  function searchServer(event) {
    event.preventDefault();
    const data = new FormData(event.target);
    const server_ID = data.get("serverID");
    console.log("Server ID:", server_ID);
  }

  function handleAddChannel(event) {
    event.preventDefault();
    if (newChannelName.trim() && selectedServer) {
      
      const updatedServer = {
        ...selectedServer,
        channels: [...(selectedServer.channels || []), newChannelName]
      };
      
      
      setSelectedServer(updatedServer);
      
      setNewChannelName("");
      setShowAddChannelForm(false);
    }
  }

  return (
    <div style={{ display: "flex" }}>
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
        <Search id="serverSearchInput" return={handleServerSearch} />
        <List sx={{display: "flex", flexDirection: "column", width: "100%", height: "100%", overflow: "auto"}}>
          <div id="serverBadgeHolder">
            {servers
              .filter((server) => server.name.toLowerCase().includes(query.toLowerCase()))
              .map((server) => (
                <li
                  key={server.id}
                  onClick={() => handleServerClick(server)}
                  style={{ cursor: "pointer" }}
                >
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
      
      {selectedServer && (
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
          
          <List sx={{ overflow: 'auto', flexGrow: 1 }}>
            {selectedServer.channels && selectedServer.channels.map((channel, index) => (
              <li key={`${selectedServer.id}-channel-${index}`}>
                <ServerBadge server={{ id: `${selectedServer.id}-channel-${index}`, name: channel }} />
              </li>
            ))}
          </List>
          
          {showAddChannelForm ? (
            <div className="addChannelForm">
              <form onSubmit={handleAddChannel}>
                <input
                  type="text"
                  value={newChannelName}
                  onChange={(e) => setNewChannelName(e.target.value)}
                  placeholder="Channel name"
                />
                <div className="formButtons">
                  <button type="submit">Add</button>
                  <button type="button" onClick={() => setShowAddChannelForm(false)}>Cancel</button>
                </div>
              </form>
            </div>
          ) : (
            <div className="addChannelButton">
              <button onClick={() => setShowAddChannelForm(true)}>+ Add Channel</button>
            </div>
          )}
        </Paper>
      )}
    </div>
  );
}

export default ServerList;