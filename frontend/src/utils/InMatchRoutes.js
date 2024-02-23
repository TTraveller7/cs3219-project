import {useEffect, useState} from "react";
import {Navigate, Outlet, useOutletContext} from "react-router-dom";
import {getUserInMatch} from "../model/Match";

function InMatchRoutes() {
    const [match, setMatch] = useState();
    const [isInMatch, setIsInMatch] = useState();
    const [username] = useOutletContext();

    async function checkIfInMatch() {
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
        setMatch(match)
    }

    useEffect(() => {checkIfInMatch().then()},[]);

    if (isInMatch !== undefined) {
        return isInMatch ?  <Outlet context={[username, match]}/> : <Navigate to={"/difficulty"}/>;
    }
}



export default InMatchRoutes;