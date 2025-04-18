import { Button, Paper } from "@mui/material";
import Search from "../../CommonComponents/Search/Search.jsx";
import { ArrowRight } from "@mui/icons-material";
import { useState } from "react";
import {useUser} from "../../UserContext.jsx";

function CreateGroup() {

    const serverContext = useUser();

    const [groupName, setGroupName] = useState("");

    const handleClick = () => {

        if(groupName === "")
            return;

        console.log("Creating group" + groupName);

        serverContext.CreateGroup({groupName});
    };

    return (
        <Paper sx={{ display: "flex", flexDirection: "row", overflow: "hidden" }}>
            <Search label="Group" onChange={setGroupName} value={groupName} />
            <Button onClick={handleClick}>
                <ArrowRight />
            </Button>
        </Paper>
    );
}

export default CreateGroup;
