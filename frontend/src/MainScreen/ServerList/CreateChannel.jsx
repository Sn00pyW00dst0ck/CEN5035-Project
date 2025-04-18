import { Button, Paper } from "@mui/material";
import Search from "../../CommonComponents/Search/Search.jsx";
import { ArrowRight } from "@mui/icons-material";
import { useState } from "react";
import {useUser} from "../../UserContext.jsx";

function CreateChannel() {

    const serverContext = useUser();

    const [channelName, setChannelName] = useState("");

    const handleClick = () => {

        if(channelName === "")
            return;

        console.log("Creating channel " + channelName);

        serverContext.CreateChannel({channelName});
    };

    return (
        <Paper sx={{ display: "flex", flexDirection: "row", overflow: "hidden" }}>
            <Search label="Channel" onChange={setChannelName} value={channelName} />
            <Button onClick={handleClick}>
                <ArrowRight />
            </Button>
        </Paper>
    );
}

export default CreateChannel;