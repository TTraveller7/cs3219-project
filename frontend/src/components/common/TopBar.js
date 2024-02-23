import AccountCircle from '@mui/icons-material/AccountCircle';
import {
    Alert, AppBar, Box, Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle,
    IconButton, Menu, MenuItem, TextField, Toolbar, Typography
} from "@mui/material";
import React, { useState } from "react"
import axios from "axios";
import { URL_USER_SVC_LOGOUT, URL_USER_SVC_USER_DELETE } from "../../configs";
import { STATUS_CODE_DELETE, STATUS_CODE_LOGOUT, STATUS_CODE_NOT_LOGGED_IN } from "../../constants";

import { useNavigate, useOutletContext } from "react-router-dom";

function TopBar() {
    const [username] = useOutletContext();
    const [deletedUsername, setDeleteUsername] = useState('');
    const [anchorEl, setAnchorEl] = React.useState(null);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [isDeleteSuccessDialogOpen, setIsDeleteSuccessDialogOpen] = useState(false);
    const [deleteErrorMessage, setIsDeleteErrorMessage] = useState('');
    const navigate = useNavigate();

    const redirectToChangePasswordPage = () => {
        navigate('/change-password');
    }

    const logoutUser = async () => {
        const res = await axios.post(URL_USER_SVC_LOGOUT, { username }, { withCredentials: true })
            .catch((err) => {
                navigate('/');
            });
        if (res && res.status === STATUS_CODE_LOGOUT) {
            navigate('/');
        }
    }

    const handleMenu = (event) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const openDeleteDialog = () => {
        handleClose();
        setDeleteUsername('')
        setIsDeleteDialogOpen(true);
    }

    const closeDeleteDialog = () => {
        setIsDeleteDialogOpen(false);
    }

    const deleteUser = async () => {
        const res = await axios.post(
            URL_USER_SVC_USER_DELETE, { "username": deletedUsername }, { withCredentials: true })
            .catch((err) => {
                setIsDeleteErrorMessage(err.response.data.error)
            })
        if (res && res.status === STATUS_CODE_DELETE) {
            closeDeleteDialog();
            openDeleteSuccessDialog();
        }
    }

    const openDeleteSuccessDialog = () => {
        setIsDeleteSuccessDialogOpen(true)
    }

    const closeDeleteSuccessDialog = () => {
        setIsDeleteSuccessDialogOpen(false)
    }

    // have 3 difficulty levels: Easy, Medium, Hard then redirect to the respective page
    const deleteUserDialog = (
        <Dialog
            open={isDeleteDialogOpen}
            onClose={closeDeleteDialog}
        >
            <DialogTitle>We're sorry to see you go!</DialogTitle>
            <DialogContent>
                <DialogContentText>Are you sure you want to delete your account?</DialogContentText>
                <DialogContentText>Note that this action cannot be undone.</DialogContentText>
                <DialogContentText><b>To delete, please enter "{username}" into the text box below.</b></DialogContentText>
                {deleteErrorMessage && <Alert severity="error">{deleteErrorMessage}</Alert>}
                <TextField
                    variant="outlined"
                    value={deletedUsername}
                    onChange={(e) => setDeleteUsername(e.target.value)}
                    sx={{ marginTop: "1.5rem" }}
                    autoFocus
                />
            </DialogContent>
            <DialogActions>
                <Button onClick={deleteUser}>Delete</Button>
            </DialogActions>
        </Dialog>
    );

    const deleteSuccessDialogue = (
        <Dialog
            open={isDeleteSuccessDialogOpen}
            onClose={closeDeleteSuccessDialog}
        >
            <DialogTitle>Account has been successfully deleted!</DialogTitle>
            <DialogActions>
                <Button onClick={logoutUser}>Logout</Button>
            </DialogActions>
        </Dialog>
    );

    return (
        <Box sx={{flexGrow: 1}}>
            <AppBar position="static">
                <Toolbar>
                    <Typography variant="h6" component="div" sx={{flexGrow: 1}}>
                        Hi {username}!
                    </Typography>
                    <div>
                        <IconButton
                            size="large"
                            aria-label="account of current user"
                            aria-controls="menu-appbar"
                            aria-haspopup="true"
                            onClick={handleMenu}
                            color="inherit"
                        >
                            <AccountCircle/>
                        </IconButton>
                        <Menu
                            id="menu-appbar"
                            anchorEl={anchorEl}
                            anchorOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            keepMounted
                            transformOrigin={{
                                vertical: 'top',
                                horizontal: 'right',
                            }}
                            open={Boolean(anchorEl)}
                            onClose={handleClose}
                        >
                            <MenuItem onClick={redirectToChangePasswordPage}>Change Password</MenuItem>
                            <MenuItem onClick={openDeleteDialog}>Delete Account</MenuItem>
                            <MenuItem onClick={logoutUser}>Logout</MenuItem>
                        </Menu>
                    </div>
                </Toolbar>
            </AppBar>
            {deleteUserDialog}
            {deleteSuccessDialogue}
        </Box>
    )
}

export default TopBar;