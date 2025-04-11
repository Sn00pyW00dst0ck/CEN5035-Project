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
    const [activeChannel, setActiveChannel] = useState([]);
    const [messages, setMessages] = useState([]);

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
            console.log("Full JSON:", jsonData);
            setGroupList(jsonData); // ✅ Update state with response data
        } catch (error) {
            console.error("Error:", error.message);
        }
    }

    async function FetchChannels() {

        if(activeGroup == null){
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
            console.log("Full channel list JSON:", jsonData);
            setChannelList(jsonData); // ✅ Update state with response data
        } catch (error) {
            console.error("Error:", error.message);
        }
    }

    async function CreateChannel({channelName}) {

        if(group == null){
            return;
        }

        try {
            const response = await fetch("http://localhost:3000/v1/api/group/" + group.id + "/channel/", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({name:{channelName}, id:"509133f2-fd45-4f1b-8bf2-72e4cb83c018"}),
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
            console.log("Full channel list JSON:", jsonData);
            setMessages(jsonData); // ✅ Update state with response data
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
        <UserContext.Provider value={{ user, setUser, groupList, activeGroup, channelList, setActiveGroup, activeChannel, setActiveChannel, messages, FetchMessages }}>
            {children}
        </UserContext.Provider>
    );
};

// Custom hook to access user context
export const useUser = () => {
    return useContext(UserContext);
};