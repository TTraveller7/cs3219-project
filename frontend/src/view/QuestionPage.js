import * as React from 'react';
import {Box, Grid} from "@mui/material";
import {AnswerBox, QuestionBox} from '../components/QuestionPage';
import ChatBoxSection from "../components/QuestionPage/ChatBox";
import QuestionTopBar from "../components/QuestionPage/QuestionTopBar";
import {io} from "socket.io-client";
import {
    BACKEND_URL,
    URL_QUESTION_SVC_GET_CURRENT_QUESTION,
} from "../configs";
import {
    createSuccess,
    getDocument,
    loadDocument, receiveChanges,
    socketCollabServiceConfig,
    socketMatchingServiceConfig
} from "../utils/SocketClientIo";
import {useEffect, useState} from "react";
import axios from "axios";
import {createQuestion, Question} from "../model/Question";
import {useOutletContext} from "react-router-dom";

let socketMatchingServiceClient;
let socketCollabServiceClient;

function QuestionPage() {
    const [username] = useOutletContext();
    const [question, setQuestion] = useState();
    const [answer, setAnswer] = useState(`// type your answer here!`);

    async function getCurrentQuestion() {
        const res = await axios.get(URL_QUESTION_SVC_GET_CURRENT_QUESTION, {withCredentials: true})
        if (res && res.data.question) {
            setQuestion(createQuestion(res.data.question));
        }
    }

    function getCurDocument(document) {
        console.log(document)
    }

    useEffect(() => {
        socketMatchingServiceClient = io(BACKEND_URL, socketMatchingServiceConfig)
        socketCollabServiceClient = io(BACKEND_URL, socketCollabServiceConfig)
        socketMatchingServiceClient.on('fetchSuccess', (res) => {
            setQuestion(new Question(res.ID, res.name, res.description))
        })
        socketCollabServiceClient.emit(getDocument, username)
        socketCollabServiceClient.on(loadDocument, (data) => {
            setAnswer(data);
        })
        socketCollabServiceClient.on(receiveChanges, (data) => {
            setAnswer(data.delta);
        })
        getCurrentQuestion()
    }, [])

    if (question) {
        return (
            <Box>
                <QuestionTopBar socketIoMatchingClient={socketMatchingServiceClient}/>
                <Box paddingTop={2}>
                    <Grid container spacing={2}>
                        <Grid item xs={12} md={6}>
                            <QuestionBox
                                socketMatchingServiceClient={socketMatchingServiceClient}
                                socketCollabServiceClient={socketCollabServiceClient}
                                question={question}
                                answer={answer}
                            />
                        </Grid>
                        <Grid item xs={12} md={6} pr={2}>
                            <AnswerBox
                                answer={answer}
                                setAnswer={setAnswer}
                                socketCollabServiceClient={socketCollabServiceClient}
                            />
                        </Grid>
                    </Grid>
                    <ChatBoxSection/>
                </Box>
            </Box>
        )
    }
}

export default QuestionPage;
