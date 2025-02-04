import './UserBadge.css'
import {Avatar, Badge, Paper, Stack} from "@mui/material";
function UserBadge() {

    return (
        <Paper elevation={75} style={{display: "block", height: "fit-content", width: "inherit", margin: ".5rem"}} className="UserBadgeContainer" sx={{
            borderRadius: 10
        }}>
            <Stack direction = "row" spacing ={2}>
                <Badge style={{margin: "2px"}} color="primary" badgecontent=" " variant ="dot">
                    <Avatar alt = "test"></Avatar>
                </Badge>

                <div>
                    <h3>Username</h3>

                    <h6>Status</h6>
                </div>
            </Stack>
        </Paper>
    )

}

export default UserBadge;