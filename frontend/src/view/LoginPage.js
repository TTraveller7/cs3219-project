import {
    Alert,
    Box,
    Button,
    Link,
    TextField,
    Typography
} from "@mui/material";
import {useState} from "react";
import axios from "axios";
import {URL_USER_SVC_LOGIN} from "../configs";
import {STATUS_CODE_LOGIN} from "../constants";
import React from "react"
import { FormLayout } from "../components/common";
import { useNavigate } from "react-router-dom";

function SignupPage() {

    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [errorMsg, setErrorMsg] = useState("")
    const navigate = useNavigate();

    const handleLogin = async () => {
        const res = await axios.post(URL_USER_SVC_LOGIN, { username, password }, {withCredentials: true})
            .catch((err) => {
                setErrorAlert(err.response.data.error)
            })
        if (res && res.status === STATUS_CODE_LOGIN) {
            resetErrorAlert();
            navigate('/difficulty')
        }
    }

    const setErrorAlert = (msg) => {
        setErrorMsg(msg)
    }

    const resetErrorAlert = () => {
        setErrorMsg('')
    }

    return (
    <React.Fragment>
        <FormLayout>
            <Typography
                variant={"h5"}
                marginBottom={"4rem"}
                textAlign={"center"}
                fontWeight={700}
            >
                Sign in
            </Typography>
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
            <Box
                display={"flex"}
                marginTop={"3rem"}
            >
                <Button
                    variant={"contained"}
                    onClick={handleLogin}
                    sx={{ flexGrow: 1, fontWeight: 700, textTransform: "none" }}
                >
                    Sign in
                </Button>
            </Box>
        </FormLayout>
        <Box
            width={"30%"}
            margin={"0 auto"}
        >
            <Box
                display={"flex"}
                marginTop={"1rem"}
                alignItems={"center"}
                justifyContent={"center"}
                width="100%"
                variant={"subtitle2"}
            >
                <Typography variant={"subtitle2"}>Don't have an account?</Typography>
                &nbsp;
                <Link href="/signup" variant={"subtitle2"}>Sign up now!</Link>
            </Box> 
        </Box>
    </React.Fragment>
    )
}

export default SignupPage;
