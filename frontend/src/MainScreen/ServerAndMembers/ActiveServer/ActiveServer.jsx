import { Paper } from "@mui/material";
import "./ActiveServer.css";

function ActiveServer() {
    return (
        <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>

            <div className="text">
                <input type="text" placeholder="Text Message" />
                <button className="sendButton">Send</button>
            </div>

        </div>
    );
}

export default ActiveServer;
