import {Paper} from "@mui/material";
import ActiveServer from "./ActiveServer/ActiveServer.jsx";
import Members from "./Members/Members.jsx";
import {useState} from "react";

function ServerAndMembers() {

    const [visible, setVisible] = useState(false);

    return(
            <Paper elevation={3} sx={{
                borderRadius: 7.5
            }} style={{ display: "flex", margin: "1rem", width: "100%", overflow: "hidden"}}>

                <ActiveServer setVisible={setVisible}/>

                {visible && <Members/>}

            </Paper>
    )
}

export default ServerAndMembers;