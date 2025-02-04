import UserBadge from "../../UserBadge/UserBadge.jsx";
import {Paper} from "@mui/material";

function ServerList() {
    return(
            <Paper elevation={3} sx={{
                borderRadius: 7.5
            }} style={{ display: "flex", width: "25vw", height: "calc(100vh - 2rem)", margin: "1rem"}}>

                <UserBadge/>

            </Paper>
    )
}

export default ServerList;