import axios from "axios";
import {URL_USER_SVC_USER_CHANGE_PASSWORD} from "../configs";
import {STATUS_CODE_CHANGE_PASSWORD} from "../constants";
import React, {useState} from "react"
import {useNavigate} from "react-router-dom";
import {FormLayout} from "../components/common";
import {
    Alert,
    Box,
    Button, Dialog, DialogActions,
    DialogTitle,
    TextField,
    Typography
} from "@mui/material";

function ChangePassword() {
    const [oldPassword, setOldPassword] = useState('');
    const [newPassword, setNewPassword] = useState('');
    const [errorMsg, setErrorMsg] = useState('');
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const navigate = useNavigate();

    const changePassword = async () => {
        const res = await axios.post(
            URL_USER_SVC_USER_CHANGE_PASSWORD,
            {oldPassword, newPassword},
            {
                withCredentials: true,

            })
            .catch((err) => {
                setErrorAlert(err.response.data.error)
            });
        if (res && res.status === STATUS_CODE_CHANGE_PASSWORD) {
            restartErrorAlert();
            setIsDialogOpen(true);
        }
    }

    const setErrorAlert = (msg) => {
        setErrorMsg(msg)
    }

    const restartErrorAlert = () => {
        setErrorMsg('')
    }

    const closeDialog = () => {
        setIsDialogOpen(false);
        navigate('/difficulty');
    }

    const successfulPasswordChangeDialogue = (
        <Dialog
        open={isDialogOpen}
        onClose={closeDialog}
        >
            <DialogTitle>Password change successfully</DialogTitle>
            <DialogActions>
                <Button onClick={closeDialog}>Login</Button>
            </DialogActions>
        </Dialog>
    );

    return (
        <FormLayout>
            <Typography variant={"h5"} marginBottom={"4rem"} textAlign={"center"} fontWeight={700}>Change
                password</Typography>
            <TextField
                label="Old password"
                variant="outlined"
                type={"password"}
                value={oldPassword}
                onChange={(e) => setOldPassword(e.target.value)}
                sx={{marginBottom: "1.5rem"}}
                autoFocus
            />
            <TextField
                label="Password"
                variant="outlined"
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                sx={{marginBottom: "2rem"}}
            />
            {errorMsg !== '' && <Alert severity="error">{errorMsg}</Alert>}
            <Box display={"flex"} marginTop={"3rem"}>
                <Button
                    onClick={changePassword}
                    variant={"contained"}
                    sx={{
                        flexGrow: 1,
                        fontWeight: 700,
                        textTransform: "none"
                    }}
                >
                    Change password
                </Button>
            </Box>
            {successfulPasswordChangeDialogue}
        </FormLayout>
    )
}

export default ChangePassword;