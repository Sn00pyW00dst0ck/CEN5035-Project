import React, {useState, useCallback} from 'react';
import {
    List,
    Paper,
    TextField,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    Avatar,
    Select,
    MenuItem
} from "@mui/material";
import ServerBadge from "./ServerBadge/ServerBadge.jsx";
import Search from "../../CommonComponents/Search/Search.jsx";
import "./ServerList.css";
import {useUser} from "../../UserContext.jsx";
import CreateGroup from "./CreateGroup.jsx";
import CreateChannel from "./CreateChannel.jsx";

// Constants for user statuses
const USER_STATUSES = [
    {value: 'online', label: 'Online', color: 'green'},
    {value: 'away', label: 'Away', color: 'orange'},
    {value: 'do-not-disturb', label: 'DND', color: 'red'},
    {value: 'invisible', label: 'Invisible', color: 'gray'}
];


// Create a context for global server state
//export const ServerContext = React.createContext();

function ProfileEditModal({
                              open,
                              onClose,
                              user,
                              onUpdateUser
                          }) {
    const [editedUser, setEditedUser] = useState({...user});
    const [selectedImage, setSelectedImage] = useState(null);

    const handleImageChange = (event) => {
        const file = event.target.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setSelectedImage(reader.result);
                setEditedUser(prev => ({...prev, icon: reader.result}));
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
                        style={{display: 'none'}}
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
        {value: 'online', label: 'Online', color: 'green'};

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
                <div style={{fontWeight: 'bold'}}>{user}</div>
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

function ServerList({onChannelSelect}) {

    const serverContext = useUser();

    const [query, setQuery] = useState('');

    const handleServerSearch = useCallback((value) => {
        setQuery(value);
    }, []);


    function handleServerClick(server){
        console.log("Active server changing to " + server.id);
        serverContext.setActiveGroup(server);
    }

    const filteredServers = serverContext.groupList == null ? [] : serverContext.groupList.filter((server) =>
        server.name.toLowerCase().includes(query.toLowerCase())
    );

    const searchServer = useCallback((event) => {
        event.preventDefault();
    }, [/*state.joinServerInput*/]);

    const handleUpdateUser = useCallback((updatedUser) => {
        setState(prev => ({
            ...prev,
            YourUser: updatedUser,
            isProfileEditOpen: false
        }));
    }, []);

    return (
        <div>
            <div style={{display: "flex"}}>
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
                        user={serverContext.user.username}
                        status="no status"
                        online={serverContext.user.online}
                        img={serverContext.user.icon}
                        about={serverContext.user.status}
                        onEditProfile={() => handleInputChange('isProfileEditOpen', true)}
                    />
                    <Search
                        id="serverSearchInput"
                        onChange={handleServerSearch}
                    />
                    <List sx={{
                        display: "flex",
                        flexDirection: "column",
                        width: "100%",
                        height: "100%",
                        overflow: "auto"
                    }}>
                        <div id="serverBadgeHolder">
                            {filteredServers.map((server) => (
                                <li
                                    key={server.id}
                                    style={{cursor: "pointer"}}
                                >
                                    <ServerBadge server={server} onClickIn={handleServerClick}/>
                                </li>
                            ))}
                        </div>
                    </List>
                    <div>
                        <CreateGroup/>
                    </div>
                </Paper>

                {serverContext.id !== '' && (
                    <Paper elevation={3} sx={{
                        borderRadius: 7.5,
                        display: "flex",
                        flexDirection: "column",
                        width: "15rem",
                        height: "calc(100vh - 2rem)",
                        margin: "0rem",
                        marginTop: "1rem",
                        overflow: "hidden"
                    }}>
                        <h3 style={{textAlign: "center", marginTop: "1rem"}}>
                            {serverContext.activeGroup.name} Channels
                        </h3>

                        <List sx={{overflow: 'auto', flexGrow: 1}}>
                            {serverContext.channelList && serverContext.channelList.map((channel, index) => (

                                <li
                                    key={`${serverContext.activeGroup.id}-channel-${index}`}
                                    style={{
                                        cursor: "pointer",
                                        backgroundColor: /*state.selectedChannel === channel ?*/ 'rgba(25, 118, 210, 0.08)' /*: 'transparent'*/
                                    }}
                                    onClick={() => serverContext.setActiveChannel(channel)}
                                >
                                    <ServerBadge
                                        server={{
                                            id: `${channel.id}-channel-${index}`,
                                            name: channel.name
                                        }}
                                    />
                                </li>
                            ))}
                        </List>

                        <CreateChannel/>

                    </Paper>
                )}
            </div>

            <ProfileEditModal
                open={false/*state.isProfileEditOpen*/}
                onClose={() => handleInputChange('isProfileEditOpen', false)}
                user={"online"/*state.YourUser*/}
                onUpdateUser={handleUpdateUser}
            />
        </div>

    );
}

export default ServerList;