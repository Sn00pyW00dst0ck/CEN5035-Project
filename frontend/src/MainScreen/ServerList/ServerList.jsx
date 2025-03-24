import React, { useState } from 'react';
import UserBadge from "../../UserBadge/UserBadge.jsx";
import { List, Paper, TextField, Dialog, DialogTitle, DialogContent, DialogActions, Button, Avatar, Select, MenuItem } from "@mui/material";
import ServerBadge from "./ServerBadge/ServerBadge.jsx";
import Search from "../../CommonComponents/Search/Search.jsx";
import "./ServerList.css";

const USER_STATUSES = [
  { value: 'online', label: 'Online', color: 'green' },
  { value: 'away', label: 'Away', color: 'orange' },
  { value: 'do-not-disturb', label: 'DND', color: 'red' },
  { value: 'invisible', label: 'Invisible', color: 'gray' }
];

function ProfileEditModal({ 
  open, 
  onClose, 
  user, 
  onUpdateUser 
}) {
  const [editedUser, setEditedUser] = useState({ ...user });
  const [selectedImage, setSelectedImage] = useState(null);

  const handleImageChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setSelectedImage(reader.result);
        setEditedUser(prev => ({ ...prev, icon: reader.result }));
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSave = () => {
    onUpdateUser(editedUser);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Edit Profile</DialogTitle>
      <DialogContent>
        <div style={{ 
          display: 'flex', 
          flexDirection: 'column', 
          alignItems: 'center', 
          gap: '1rem' 
        }}>
          <input
            accept="image/*"
            style={{ display: 'none' }}
            id="profile-image-upload"
            type="file"
            onChange={handleImageChange}
          />
          <label htmlFor="profile-image-upload">
            <Avatar 
              src={selectedImage || editedUser.icon} 
              sx={{ 
                width: 100, 
                height: 100, 
                cursor: 'pointer' 
              }} 
            />
          </label>
          
          <TextField
            fullWidth
            label="Username"
            value={editedUser.name}
            onChange={(e) => setEditedUser(prev => ({ 
              ...prev, 
              name: e.target.value 
            }))}
          />
          
          <Select
            fullWidth
            value={editedUser.status}
            onChange={(e) => setEditedUser(prev => ({ 
              ...prev, 
              status: e.target.value 
            }))}
          >
            {USER_STATUSES.map((status) => (
              <MenuItem key={status.value} value={status.value}>
                {status.label}
              </MenuItem>
            ))}
          </Select>
          
          <TextField
            fullWidth
            label="About Me"
            multiline
            rows={4}
            value={editedUser.about || ''}
            onChange={(e) => setEditedUser(prev => ({ 
              ...prev, 
              about: e.target.value 
            }))}
          />
        </div>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSave} color="primary">Save</Button>
      </DialogActions>
    </Dialog>
  );
}

function CustomUserBadge({ 
  user, 
  status, 
  online, 
  img, 
  about,
  onEditProfile 
}) {
  const statusConfig = USER_STATUSES.find(s => s.value === status) || 
                       { value: 'online', label: 'Online', color: 'green' };

  return (
    <div style={{ 
      display: 'flex', 
      alignItems: 'center', 
      padding: '1rem', 
      position: 'relative' 
    }}>
      <Avatar 
        src={img} 
        alt={user} 
        sx={{ 
          width: 50, 
          height: 50, 
          marginRight: '1rem' 
        }} 
      />
      <div>
        <div style={{ fontWeight: 'bold' }}>{user}</div>
        <div 
          style={{ 
            color: statusConfig.color, 
            display: 'flex', 
            alignItems: 'center' 
          }}
        >
          <div 
            style={{ 
              width: 10, 
              height: 10, 
              borderRadius: '50%', 
              backgroundColor: statusConfig.color, 
              marginRight: '0.5rem' 
            }}
          />
          {statusConfig.label}
        </div>
        {about && (
          <div style={{ 
            color: 'gray', 
            fontSize: '0.8rem', 
            marginTop: '0.25rem' 
          }}>
            {about}
          </div>
        )}
      </div>
      <Button 
        onClick={onEditProfile}
        style={{ 
          position: 'absolute', 
          right: 10, 
          top: '50%', 
          transform: 'translateY(-50%)' 
        }}
      >
        Edit
      </Button>
    </div>
  );
}

function ServerList({ servers }) {
  const [query, setQuery] = useState("");
  const [selectedServer, setSelectedServer] = useState(null);
  const [newChannelName, setNewChannelName] = useState("");
  const [showAddChannelForm, setShowAddChannelForm] = useState(false);
  const [isProfileEditOpen, setIsProfileEditOpen] = useState(false);

  // Updated YourUser state to be more dynamic
  const [YourUser, setYourUser] = useState({
    name: "Your Username",
    status: "online",
    online: true,
    icon: "/default-profile.png", // Default profile image
    about: "Hello! I'm using the app."
  });

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

  const handleUpdateUser = (updatedUser) => {
    setYourUser(updatedUser);
  };

  return (
    <>
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
          <CustomUserBadge 
            user={YourUser.name} 
            status={YourUser.status} 
            online={YourUser.online} 
            img={YourUser.icon}
            about={YourUser.about}
            onEditProfile={() => setIsProfileEditOpen(true)}
          />
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

      <ProfileEditModal 
        open={isProfileEditOpen}
        onClose={() => setIsProfileEditOpen(false)}
        user={YourUser}
        onUpdateUser={handleUpdateUser}
      />
    </>
  );
}

export default ServerList;