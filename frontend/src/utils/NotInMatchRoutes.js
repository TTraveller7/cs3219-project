import {useEffect, useState} from "react";
import {Navigate, Outlet, useOutletContext} from "react-router-dom";
import {getUserInMatch} from "../model/Match";

function NotInMatchRoutes() {
    const [isInMatch, setIsInMatch] = useState();
    const [username] = useOutletContext();

    const checkIfInMatch = async () => {
        const match = await getUserInMatch(username);
        if (match === undefined) {
            setIsInMatch(false)
            return;
        }
        if (match.isEnded) {
            setIsInMatch(false)
            return;
        }
        setIsInMatch(true)
    }

    useEffect(() => {
        checkIfInMatch()
    },[]);

    if (isInMatch !== undefined) {
        return isInMatch ? <Navigate to={"/question"}/> : <Outlet context={[username]}/>;
    }
}

export default NotInMatchRoutes;