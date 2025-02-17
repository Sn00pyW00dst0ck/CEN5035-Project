import {Button, Paper, TextField} from "@mui/material";
import "./ActiveServer.css";

function ActiveServer({setVisible}) {
    return (
        <div style={{ display: "flex", flexDirection: "column", height: "100%", width: "100%" }}>

            <Button onClick={() => setVisible((prev) => !prev)} sx={{
                width: "4rem",
                height: "4rem",
                margin: "auto",
                marginRight: "0",
                marginTop: "0",
                borderRadius: 7.5,
                borderTopRightRadius: 0,
                borderBottomRightRadius: 0
            }}>Show/hide</Button>

            <Paper elevation={7} sx={{borderRadius: 7.5,
                borderTopLeftRadius: 0,
                borderTopRightRadius: 0,
                borderBottomRightRadius: 0}} className="text">

                <TextField sx={{
                    width: "100%",
                    margin: ".5rem",
                    marginLeft: "1rem",
                }} placeholder="Text Message" />

                <button className="sendButton">Send</button>

            </Paper>


        </div>
    );
}

export default ActiveServer;
