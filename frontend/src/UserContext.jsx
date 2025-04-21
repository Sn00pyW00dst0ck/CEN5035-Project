// src/UserContext.js
import React, {createContext, useContext, useEffect, useState} from 'react';

// Create the context to hold user data
const UserContext = createContext();

// Provider component to wrap the part of your app that needs access to user data
export const UserProvider = ({ children }) => {
    const [user, setUser] = useState({created_at:'', id: '', profile_pic: '', username: ''});

    const [groupList, setGroupList] = useState([]);
    const [channelList, setChannelList] = useState([]);

    const [activeGroup, setActiveGroup] = useState({created_at:'', description:'', group:'', id:'', name:''});
    const [activeChannel, setActiveChannel] = useState({});
    const [messages, setMessages] = useState([]);

    const [userMap, setUserMap] = useState({});

    async function FetchGroups() {

        if(user === null || user.id === ''){
            return;
        }

        try {
            const response = await fetch("http://localhost:3000/v1/api/group/search", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({members:[user.id]}),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const jsonData = await response.json();
            setGroupList(jsonData); // ✅ Update state with response data
        } catch (error) {
            console.error("Error:", error.message);
        }
    }

    async function FetchChannels() {

        if(activeGroup === null || activeGroup.id === ''){
            return;
        }

        try {
            const response = await fetch("http://localhost:3000/v1/api/channel/search", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({groups:[{group:activeGroup.id}]}),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const jsonData = await response.json();
            setChannelList(jsonData); // ✅ Update state with response data
        } catch (error) {
            console.error("Error:", error.message);
        }
    }

    async function CreateGroup({groupName}) {

        try {
            const response = await fetch("http://localhost:3000/v1/api/group/", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    name:groupName,
                    id:crypto.randomUUID(),
                    description:"",
                    members:[user.id],
                }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const jsonData = await response.json();

            await FetchChannels();

        } catch (error) {
            console.error("Error:", error.message);
        }
    }

    async function CreateChannel({channelName}) {

        if(activeGroup.id == null)
            return;

        try {
            const response = await fetch("http://localhost:3000/v1/api/group/" + activeGroup.id + "/channel/", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    name:channelName,
                    id:crypto.randomUUID(),
                    group:activeGroup.id
                    }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const jsonData = await response.json();

            await FetchChannels();

        } catch (error) {
            console.error("Error:", error.message);
        }
    }

    async function FetchMessages() {

        if(activeGroup == null){
            return;
        }

        if(activeChannel == null || activeChannel.id === '' || activeChannel.id ==='0'){
            return;
        }
        try {
            const response = await fetch("http://localhost:3000/v1/api/message/search", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({channel:[activeChannel.id]}),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const jsonData = await response.json();

            jsonData.sort((a, b) => new Date(a.created_at) - new Date(b.created_at));

            setMessages(jsonData); // ✅ Update state with response data
        } catch (error) {
            setMessages([]);
            console.error("Error:", error.message);
        }
    }

    async function FetchUserFromId({userID}){
        try {
            const response = await fetch("http://localhost:3000/v1/api/account/search", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    id: {userID}
                }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const jsonData = await response.json();

            return jsonData[0].name;

        } catch (error) {
            console.error("Error:", error.message);
        }
    }
    useEffect(() => {
        if (activeGroup == null) return;

        const FetchMessagesAndUsers = async () => {

            try {

                // Extract unique user IDs
                const uniqueUserIds = [...new Set(messages.map(msg => msg.author))];

                // Fetch all user names
                const userEntries = await Promise.all(
                    uniqueUserIds.map(async (userID) => {
                        const res = await fetch("http://localhost:3000/v1/api/account/search", {
                            method: "POST",
                            headers: { "Content-Type": "application/json" },
                            body: JSON.stringify({ id: [userID] }),
                        });

                        if (!res.ok) throw new Error(`User fetch failed for ${userID}`);

                        const userData = await res.json();

                        return [userID, userData[0]?.username ?? "Unknown"];
                    })
                );

                const userMapObj = Object.fromEntries(userEntries);
                setUserMap(userMapObj);
            } catch (error) {
                console.error("Error:", error.message);
            }
        };

        FetchMessagesAndUsers();
    }, [activeGroup, activeChannel]);


    async function SendMessage(message){

        if(activeGroup.id == null)
            return;

        if(activeChannel.id == null)
            return;

        try {
            const response = await fetch("http://localhost:3000/v1/api/group/" + activeGroup.id + "/channel/" + activeChannel.id + "/message", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    author:user.id,
                    body: message,
                    channel:activeChannel.id,
                    id:crypto.randomUUID(),
                    pinned:false
                }),
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            const jsonData = await response.json();

            await FetchChannels();

        } catch (error) {
            console.error("Error:", error.message);
        }
    }

    useEffect(() => {

        FetchChannels();

    }, [activeGroup])

    useEffect(() => {

        FetchGroups();

    }, [user])

    useEffect(() => {

        FetchMessages();

    }, [activeChannel])

    return (
        <UserContext.Provider value={{ user, userMap, setUser, groupList, activeGroup, channelList, setActiveGroup, activeChannel, setActiveChannel, messages, CreateGroup, CreateChannel, SendMessage, FetchUserFromId}}>
            {children}
        </UserContext.Provider>
    );
};

// Custom hook to access user context
export const useUser = () => {
    return useContext(UserContext);
};