import { Box } from "@mui/material";
import React from "react"

const Form = ({ children }) => {
    return (
        <Box display={"flex"} flexDirection={"column"} padding={"4rem"}>
            <Box width={"30%"} margin={"0 auto"}>
                <Box display={"flex"} flexDirection={"column"} width={"100%"} boxSizing={"border-box"} sx={{
                    backgroundColor: "#FFFFFF",
                    boxShadow: "0 0 3px #ccc",
                    padding: "2rem",
                    borderRadius: 2,
                }}>
                    {children}
                </Box>
            </Box>
        </Box>
    )
}

export default Form;