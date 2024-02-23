import {
    Alert,
    Box,
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogContentText,
    DialogTitle,
    Link,
    TextField,
    Typography
} from "@mui/material";
import {useState} from "react";
import axios from "axios";
import {URL_USER_SVC_USER_CREATE} from "../configs";
import {STATUS_CODE_CREATED} from "../constants";
import React from "react"
import { FormLayout } from "../components/common";
import { useNavigate } from "react-router-dom";

function SignupPage() {
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [isDialogOpen, setIsDialogOpen] = useState(false)
    const [dialogTitle, setDialogTitle] = useState("")
    const [dialogMsg, setDialogMsg] = useState("")
    const [errorMsg, setErrorMsg] = useState('');
    const navigate = useNavigate();

    const handleSignup = async () => {
        const res = await axios.post(
            URL_USER_SVC_USER_CREATE, { username, password }, {withCredentials: true})
            .catch((err) => {
                setErrorAlert(err.response.data.error)
            })
        if (res && res.status === STATUS_CODE_CREATED) {
            restartErrorAlert()
            setSuccessDialog('Account successfully created')
        }
    }

    const closeDialog = () => setIsDialogOpen(false)

    const setSuccessDialog = (msg) => {
        setIsDialogOpen(true)
        setDialogTitle('Success')
        setDialogMsg(msg)
    }

    const setErrorAlert = (msg) => {
        setErrorMsg(msg)
    }

    const restartErrorAlert = () => {
        setErrorMsg('')
    }

    const redirectToLogin = () => {
        let path = '/login';
        navigate(path);
    }

    const dialog = (
        <Dialog
            open={isDialogOpen}
            onClose={closeDialog}
        >
            <DialogTitle>{dialogTitle}</DialogTitle>
            <DialogContent>
                <DialogContentText>{dialogMsg}</DialogContentText>
            </DialogContent>
            <DialogActions>
                <Button onClick={redirectToLogin}>Login</Button>
            </DialogActions>
        </Dialog>
    );

    return (
        <React.Fragment>
            <FormLayout>
                <Typography variant={"h5"} marginBottom={"4rem"} textAlign={"center"} fontWeight={700}>Create an
                    account</Typography>
                <TextField
                    label="Username"
                    variant="outlined"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    sx={{marginBottom: "1.5rem"}}
                    autoFocus
                />
                <TextField
                    label="Password"
                    variant="outlined"
                    type="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    sx={{marginBottom: "2rem"}}
                />
                {errorMsg !== '' && <Alert severity="error">{errorMsg}</Alert>}
                <Box display={"flex"} marginTop={"3rem"}>
                    <Button variant={"contained"} onClick={handleSignup}
                            sx={{flexGrow: 1, fontWeight: 700, textTransform: "none"}}>Sign up</Button>
                </Box>
                {dialog}
            </FormLayout>
            <Box width={"30%"} margin={"0 auto"}>
                <Box display={"flex"} marginTop={"1rem"} alignItems={"center"} justifyContent={"center"} width="100%"
                     variant={"subtitle2"}>
                    <Typography variant={"subtitle2"}>Already have an account?</Typography>
                    &nbsp;
                    <Link href="/login" variant={"subtitle2"}>Sign in</Link>
                </Box>
            </Box>
        </React.Fragment>
    )
}

export default SignupPage;
