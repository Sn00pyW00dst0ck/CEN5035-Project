import ServerBadge from "../../../ServerList/ServerBadge/ServerBadge.jsx";
import Search from "../../../../CommonComponents/Search/Search.jsx";
import {Button, Paper} from "@mui/material";
import MenuIcon from "@mui/icons-material/Menu";
import * as React from "react";

function MenuBar({ setVisible }) {
    return(
        <Paper
            sx={{
                display: "flex",
                flexDirection: "row",
                overflow: "hidden"
            }}>

            <div style={{display: "flex", flexShrink: "0", width: "15rem"}}>
                <ServerBadge server ={{id: 6, name: "Alice", icon: "public/vite.svg"}}/>
            </div>

            <Search label="Search messages"/>

            <Button
                onClick={() => setVisible((prev) => !prev)}
                sx={{
                    width: "4rem",
                    height: "4rem",
                    alignSelf: "center",
                    margin: "auto",
                    marginRight: "0",
                    borderRadius: 7.5,
                    borderTopRightRadius: 0,
                    borderBottomRightRadius: 0
                }}
            >
                <MenuIcon/>
            </Button>
        </Paper>
    );
}

export default MenuBar;