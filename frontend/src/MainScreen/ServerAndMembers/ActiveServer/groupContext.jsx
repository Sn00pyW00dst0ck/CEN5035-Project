import React, {createContext, useContext, useEffect, useState} from 'react';

// Create the context to hold user data
const GroupContext = createContext();

// Provider component to wrap the part of your app that needs access to user data
export const GroupProvider = ({ children, groupIn }) => {
    const [group, setGroup] = useState({groupIn});
    const [channelList, setChannelList] = useState([]);



    async function FetchChannels() {

        if(group == null){
            return;
        }

        try {
            const response = await fetch("http://localhost:3000/v1/api/channel/search", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({groups:[{group:group.id}]}),
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

    useEffect(() => {

        FetchChannels();

    }, [group])

    return (
        <GroupContext.Provider value={{ group, setGroup, channelList, setChannelList}}>
            {children}
        </GroupContext.Provider>
    );
};

// Custom hook to access user context
export const useGroup = () => {
    return useContext(GroupContext);
};