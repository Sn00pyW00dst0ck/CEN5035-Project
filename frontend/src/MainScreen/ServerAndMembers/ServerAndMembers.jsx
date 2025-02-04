import {Paper} from "@mui/material";
import ActiveServer from "./ActiveServer/ActiveServer.jsx";
import Members from "./Members/Members.jsx";

function ServerAndMembers() {
    return(
        <div>
            <Paper elevation={3} sx={{
                borderRadius: 7.5
            }} style={{ display: "flex", width: "calc(75vw - 4rem)", height: "calc(100vh - 2rem)", margin: "1rem"}}>


                <div style={{display: "block", width: "100%", height: "100%"}}>
                    <ActiveServer/>
                </div>

                <div style={{marginLeft: "auto", marginRight: "0"}}>
                    <Members/>
                </div>

            </Paper>
        </div>

    )
}

export default ServerAndMembers;