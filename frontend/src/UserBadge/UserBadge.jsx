import * as React from 'react';
import './UserBadge.css'
import {Avatar, Badge, Box, Icon, Paper, Stack, styled, ThemeProvider} from "@mui/material";
function UserBadge() {

    return (
        <Paper elevation={3} className="test">
            <Stack direction = "row" spacing ={2}>
                <Badge color="primary" badgecontent=" " variant ="dot">
                    <Avatar className="test" alt = "test"></Avatar>
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