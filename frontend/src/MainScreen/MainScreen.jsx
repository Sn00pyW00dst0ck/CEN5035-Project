import './MainScreen.css'
import ServerAndMembers from "./ServerAndMembers/ServerAndMembers.jsx";
import ServerList from "./ServerList/ServerList.jsx";

const servers = [
    { id: 1, name: "test1", icon: "public/vite.svg"},
    { id: 2, name: "test2", icon: "public/vite.svg"},
    { id: 3, name: "Test1", icon: "public/vite.svg"},
    { id: 4, name: "Test2", icon: "public/vite.svg"},
    { id: 5, name: "thisIsATest1", icon: "public/vite.svg"},
    { id: 6, name: "Alice", icon: "public/vite.svg"}
]

function MainScreen() {

    return (
        <div className="ColorBox" style={{ display: "flex", height: "100vh", width: "100vw"}}>
            <ServerList servers={servers} />
            <ServerAndMembers />
        </div>

    )

}

export default MainScreen;
