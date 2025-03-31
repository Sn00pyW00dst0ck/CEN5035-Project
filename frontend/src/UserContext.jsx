// src/UserContext.js
import React, {createContext, useContext, useEffect, useState} from 'react';

// Create the context to hold user data
const UserContext = createContext();



// Provider component to wrap the part of your app that needs access to user data
export const UserProvider = ({ children }) => {
    const [user, setUser] = useState({created_at:'', id: '', profile_pic: '', username: ''});

    const [groupList, setGroupList] = useState([]);

    async function FetchGroups() {
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

    useEffect(() => {

        FetchGroups();

    }, [user, setUser])

    return (
        <UserContext.Provider value={{ user, setUser, groupList}}>
            {children}
        </UserContext.Provider>
    );
};

// Custom hook to access user context
export const useUser = () => {
    return useContext(UserContext);
};