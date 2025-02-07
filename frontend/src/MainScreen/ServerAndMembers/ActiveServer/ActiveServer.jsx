import {Paper} from "@mui/material";

function ActiveServer() {
    return (
        <div style={{ display: "flex", height: "100%", width: "100%" , marginLeft: "0", marginRight: "0"}}>

            <Paper elevation={100} sx={{
                borderRadius: 7.5,
                borderTopLeftRadius: 0,
                borderTopRightRadius: 0,
                borderBottomRightRadius: 0
            }} style={{ display: "flex", height:"4rem", width: "inherit", marginBottom: "0", marginTop: "auto"}}>

                abc

            </Paper>
        </div>
    )
}

export default ActiveServer;