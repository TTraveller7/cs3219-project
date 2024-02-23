import {Button, Stack} from "@mui/material";
import * as React from "react";
import {useOutletContext} from "react-router-dom";
import axios from "axios";
import {URL_QUESTION_SVC_ANSWER_CREATE} from "../../configs";
import {createNewDocument, saveDocument} from "../../utils/SocketClientIo";

function ActionButtons({socketMatchingServiceClient, socketCollabServiceClient, question, answer}) {
    const [username, match] = useOutletContext();

    function nextQuestion() {
        socketMatchingServiceClient.emit('fetchingQuestion', username, question.id)
        socketCollabServiceClient.emit(createNewDocument, username)
    }

    async function saveQuestion() {
        await axios.post(URL_QUESTION_SVC_ANSWER_CREATE,
            {
                questionId: question.id,
                code: answer
            },
            {withCredentials: true}
        );
    }

    return(
        <Stack
            marginLeft={2}
            marginY={2}
            direction="row"
            justifyContent="space-evenly"
            alignItems="center"
            height="50px"
        >
            <Button onClick={saveQuestion}>Save</Button>
            <Button onClick={nextQuestion} variant="contained">Next</Button>
        </Stack>
    )
}

export default ActionButtons;