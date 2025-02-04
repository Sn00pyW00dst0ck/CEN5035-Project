import UserBadge from "../../../UserBadge/UserBadge.jsx";
import {Paper} from "@mui/material";

function Members() {

    return (
        <Paper elevation={25} sx={{
            borderRadius: 7.5,
            borderBottomLeftRadius: 0,
            borderTopLeftRadius: 0
        }} style={{ display: "flex", width: "20vw", height: "100%"}}>

            <div style={{display: "block", width: "calc(20vw - 1rem)"}}>

                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
                <UserBadge/>
            </div>

        </Paper>
    )
}

export default Members;