import * as React from 'react';
import { Box, Paper, styled, Typography } from "@mui/material";
import ActionButtons from "./ActionButtons";

const QuestionPaper = styled(Paper)(({ theme }) => ({
    height: 'calc(90vh - 50px)',
    border: '1px solid #ced4da',
    boxShadow: 'none',
    backgroundColor: '#fff',
    marginLeft: theme.spacing(2),
    textAlign: 'center',
    color: theme.palette.text.secondary,
}));

function QuestionBox({ socketMatchingServiceClient, question, answer, socketCollabServiceClient }) {
    return (
        <Box>
            <QuestionPaper>
                {question.id}
                <
                    Typography
                    variant="h3"
                >
                    {question.name}
                </Typography>
                <Typography
                    variant="body1"
                >
                    {question.description}
                </Typography>
            </QuestionPaper>
            <ActionButtons
                answer={answer}
                question={question}
                socketMatchingServiceClient={socketMatchingServiceClient}
                socketCollabServiceClient={socketCollabServiceClient}
            />
        </Box>
    );
}

export default QuestionBox;