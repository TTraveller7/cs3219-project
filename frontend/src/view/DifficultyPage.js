import {
    Alert,
    Box,
    Button, Dialog, DialogActions,
    DialogContent, DialogContentText,
    DialogTitle,
    TextField,
    Typography
} from "@mui/material";
import React, {useState} from "react"
import { FormLayout, TopBar } from "../components/common";
import axios from "axios";
import { URL_USER_SVC_LOGOUT, URL_USER_SVC_USER_DELETE } from "../configs";
import { STATUS_CODE_DELETE, STATUS_CODE_LOGOUT, STATUS_CODE_NOT_LOGGED_IN } from "../constants";
import { useNavigate, useOutletContext } from "react-router-dom";

function DifficultyPage() {
    const navigate = useNavigate();

    const chooseDifficulty = (chosenDifficulty) => {
        navigate('/matching', { state: chosenDifficulty });
        //navigate('/question');
    }

    // have 3 difficulty levels: Easy, Medium, Hard then redirect to the respective page
    return (
        <Box sx={{ flexGrow: 1 }}>
            <TopBar />
            <FormLayout>
                <Typography
                    variant={"h5"}
                    marginBottom={"4rem"}
                    textAlign={"center"}
                    fontWeight={700}
                >
                    Choose a difficulty level
                </Typography>
                <Button onClick={() => chooseDifficulty("easy")}>Easy</Button>
                <Button onClick={() => chooseDifficulty("medium")}>Medium</Button>
                <Button onClick={() => chooseDifficulty("hard")}>Hard</Button>
            </FormLayout>
        </Box>
    )
}

export default DifficultyPage;