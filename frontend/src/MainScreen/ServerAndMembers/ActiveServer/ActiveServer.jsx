import React, {useState, useCallback} from 'react';
import {Paper, TextField} from "@mui/material";
import "./ActiveServer.css";
import MenuBar from "./MenuBar/MenuBar.jsx";
import {useUser} from "../../../UserContext.jsx";
import Search from "../../../CommonComponents/Search/Search.jsx";

function ActiveServer({
                          setVisible/*,
  selectedServer,
  selectedChannel,
  messages,
  onChannelSelect*/
                      }) {

    const serverContext = useUser();

    const [messageInput, setMessageInput] = useState('');

    const handleUpdateMessage = useCallback((value) => {
        setMessageInput(value);
    }, []);

    function handleSendMessage() {

        console.log("Send message clicked.");
        console.log('Sending message:', messageInput);

        if (messageInput !== "" && serverContext.activeGroup && serverContext.activeChannel) {
            // In a real app, this would send to backend

            serverContext.SendMessage(messageInput);
            setMessageInput('');

        }
    }

    // Render nothing if no server selected
    if (!serverContext.activeGroup) {
        return (
            <div style={{
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                width: '100%'
            }}>
                Select a server to start chatting
            </div>
        );
    }

    return (
        <div style={{
            display: "flex",
            flexDirection: "column",
            height: "100%",
            width: "100%"
        }}>
            <MenuBar
                setVisible={setVisible}
                /*selectedServer={serverContext.group}
                selectedChannel={serverContext.channel}
                onChannelSelect={onChannelSelect}*/
            />

            {/* Message Display Area */}
            <div style={{ flexGrow: 1, overflowY: 'auto', padding: '1rem' }}>
                {serverContext.messages.length === 0 ? (
                    <p>No messages</p>
                ) : (
                    serverContext.messages.map(message => (
                        <div
                            key={message.id}
                            style={{
                                marginBottom: '0.5rem',
                                textAlign: 'left',
                                wordWrap: 'break-word',
                                whiteSpace: 'pre-wrap',
                            }}
                        >
                            <strong>{serverContext.userMap[message.author] || 'Unknown'}:</strong> {message.body}
                        </div>
                    ))
                )}
            </div>

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
                <Search onChange={handleUpdateMessage} label = "Message" value={messageInput} />
                <button
                    className="sendButton"
                    onClick={handleSendMessage}
                >
                    Send
                </button>
            </Paper>
        </div>
    );
}

export default ActiveServer;