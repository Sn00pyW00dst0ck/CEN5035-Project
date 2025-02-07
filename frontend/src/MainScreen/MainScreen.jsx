import './MainScreen.css'
import ServerAndMembers from "./ServerAndMembers/ServerAndMembers.jsx";
import ServerList from "./ServerList/ServerList.jsx";
function MainScreen() {

    return (
        <div className="ColorBox" style={{ display: "flex", height: "100vh", width: "100vw"}}>
            <ServerList />
            <ServerAndMembers />
        </div>

    )

}

export default MainScreen;
