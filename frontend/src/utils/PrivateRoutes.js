import {Outlet, Navigate} from 'react-router-dom'
import axios from "axios";
import {URL_USER_SVC_VALIDATE} from "../configs";
import {STATUS_CODE_NOT_LOGGED_IN} from "../constants";
import {useEffect, useState} from "react";

function PrivateRoutes() {
    const [username, setUsername] = useState();
    const [isLoggedIn, setIsLoggedIn] = useState();

    const verifyLogin = async () => {
        const res = await axios.post(URL_USER_SVC_VALIDATE, {}, {withCredentials: true})
            .catch((err) => {
                if (err.response.status === STATUS_CODE_NOT_LOGGED_IN) {
                    setIsLoggedIn(false);
                }
            });
        if (res && res.data.username) {
            setUsername(res.data.username);
            setIsLoggedIn(true);
        }
    }

    useEffect(() => {verifyLogin()},[]);

    if (isLoggedIn !== undefined) {
        return isLoggedIn ? <Outlet context={[username]} /> : <Navigate to={"/question"} />
    }
}

export default PrivateRoutes;