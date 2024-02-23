import {AppBar, Box, Button, Toolbar, Typography} from "@mui/material";
import {
    leaveSuccess,
    LEAVING_ROOM,
} from "../../utils/SocketClientIo";
import {useEffect} from "react";
import {useNavigate, useOutletContext} from "react-router-dom";

function QuestionTopBar({socketIoMatchingClient}) {
    const [username] = useOutletContext();
    const navigate = useNavigate();

    function leaveRoom() {
        socketIoMatchingClient.emit(LEAVING_ROOM, username)
    }

    useEffect(() => {
        socketIoMatchingClient.on(leaveSuccess, () => {
            navigate('/difficulty')
        })
    }, [])

    return (
        <Box sx={{ flexGrow: 1 }}>
            <AppBar position="static">
                <Toolbar>
                    <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                        Match
                    </Typography>
                    <Button
                        color="inherit"
                        onClick={leaveRoom}
                    >
                        Leave room
                    </Button>
                </Toolbar>
            </AppBar>
        </Box>
    )
}

export default QuestionTopBar;