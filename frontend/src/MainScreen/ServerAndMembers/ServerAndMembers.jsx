import {Paper} from "@mui/material";
import ActiveServer from "./ActiveServer/ActiveServer.jsx";
import Members from "./Members/Members.jsx";

function ServerAndMembers() {
    return(
            <Paper elevation={3} sx={{
                borderRadius: 7.5
            }} style={{ display: "flex", margin: "1rem", width: "100%"}}>

                <ActiveServer/>
                

            </Paper>
    )
}

export default ServerAndMembers;