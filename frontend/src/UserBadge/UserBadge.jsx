import './UserBadge.css'
import {Avatar, Badge, Paper, Stack} from "@mui/material";
function UserBadge(props) {

    return (
        <Paper elevation={75} style={{display: "flex", margin: ".5rem", marginBottom: "0"}} className="UserBadgeContainer" sx={{
            borderRadius: 10
        }}>
            <Stack direction = "row" spacing ={2}>
                <Badge invisible={!props.online} style={{margin: "2px"}} color="primary" badgecontent=" " variant ="dot">
                    <Avatar src ={props.img} alt = "test"></Avatar>
                </Badge>

                <div>
                    <h3>{props.user}</h3>

                    <h6>{props.status}</h6>
                </div>
            </Stack>
        </Paper>
    )

}

export default UserBadge;