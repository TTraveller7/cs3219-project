import {
    Alert,
    Box, Button, Typography
} from "@mui/material";
import { FormLayout } from "../components/common";
import {
    BACKEND_URL, URL_MATCHING_SVC_GET_USER,
    URL_MATCHING_SVC_MATCHING, URL_QUESTION_SVC_QUESTION_START
} from "../configs";
import axios from "axios";
import React, { useEffect, useState } from "react"
import { CountdownCircleTimer } from 'react-countdown-circle-timer'
import {useLocation, useNavigate, useOutletContext} from "react-router-dom";
import {socketMatchingServiceConfig,
    matchFail,
    MATCHING_TIMEOUT,
    MATCHING_USER, matchPending,
    matchSuccess,} from "../utils/SocketClientIo";
import {io} from "socket.io-client";

let socketMatchingServiceClient;

function MatchingPage() {
    const { state } = useLocation();
    const navigate = useNavigate();
    const [errorMsg, setErrorMsg] = useState("")
    const [username] = useOutletContext();
    const [key, setKey] = useState(0);
    const [reload, setReload] = useState(false);
    const [matchFound, setMatchFound] = useState(false);

    const setErrorAlert = (msg) => {
        setErrorMsg(msg)
        setReload(true);
    }

    const resetErrorAlert = () => {
        setErrorMsg('')
        setReload(false)
    }

    function restartFindingNewMatch() {
        setKey(key => key + 1);
        resetErrorAlert();
        findMatch();
        setReload(false);
    }

    const findMatch = async () => {
        await
            axios.post(URL_MATCHING_SVC_MATCHING, { "username": username, "difficulty": state }, { withCredentials: true });
    }

    function navigateBackToDifficulty() {
        leaveQueue()
        navigate('/difficulty')
    }

    function leaveQueue() {
        socketMatchingServiceClient.emit(MATCHING_TIMEOUT, username);
    }

    function findUsersInMatch() {
        // A check whehter socketMatchingServiceClient is needed is because when the timer runs this function,
        // socketMatchingServiceClient will be undefined which results in an error
        if (socketMatchingServiceClient !== undefined) {
            socketMatchingServiceClient.emit(MATCHING_USER, username);
        }
    }

    function haveMatchFunction (USER_STATUS, ANOTHER_USERNAME) {
        axios.get(URL_MATCHING_SVC_GET_USER + username, {withCredentials: true})
            .then(() => {
                setMatchFound(true);
                axios.get(URL_QUESTION_SVC_QUESTION_START, {withCredentials: true})
                    .then(() => {
                        navigate('/question')
                    })
                    .catch((err) => {
                        console.log(err)
                        setErrorAlert("Error")
                    })
            })
            .catch((err) => {
                console.log(err)
                setErrorAlert("Error")
                })
    };

    useEffect(() => {
        socketMatchingServiceClient = io(BACKEND_URL, socketMatchingServiceConfig);
        socketMatchingServiceClient.on(matchPending, (USER_STATUS) => {});
        socketMatchingServiceClient.on(matchSuccess, haveMatchFunction);
        socketMatchingServiceClient.on(matchFail, (USERNAME) => {})
        findMatch().catch(console.error);
        return () => {
            leaveQueue()
        }
    }, []) //eslint-disable-line

    const renderTime = ({ remainingTime }) => {
        if (remainingTime === 0) {
            setReload(true)
            return <div className="timer">No match found :(</div>;
        } else {
            findUsersInMatch();
        }
        return (
            <Box display="flex" flexDirection="column" alignItems="center">
                <Typography variant={"body1"}>Remaining</Typography>
                <Typography variant={"body1"}>{remainingTime}</Typography>
                <Typography variant={"body1"}>seconds</Typography>
            </Box>
        );
    };

    return (
        <React.Fragment>
            <FormLayout>
                {errorMsg !== '' && <Alert severity="error">{errorMsg}</Alert>}
                <Box width="50%" marginBottom="4rem">
                    <Button onClick={navigateBackToDifficulty}>Back</Button>
                </Box>
                <Typography
                    variant={"h5"}
                    marginBottom={"4rem"}
                    textAlign={"center"}
                    fontWeight={700}
                >
                    Matching for {state} difficulty
                </Typography>
                <Box display="flex" justifyContent="center">
                    {(!matchFound && (errorMsg === '')) && (
                        <CountdownCircleTimer
                            key={key}
                            isPlaying
                            duration={30}
                            colors={["#004777", "#F7B801", "#A30000", "#A30000"]}
                            colorsTime={[10, 6, 3, 0]}
                            onComplete={() => ({ shouldRepeat: false, delay: 1 })}
                        >
                            {renderTime}
                        </CountdownCircleTimer>
                    )}
                </Box>
                {reload && (
                    <Box marginTop={"4rem"}>
                        <Button onClick={restartFindingNewMatch}>Try to find another match</Button>
                        <Button onClick={navigateBackToDifficulty}>Choose another difficulty level</Button>
                    </Box>
                )}
            </FormLayout>
        </React.Fragment>
    )
}

export default MatchingPage;
