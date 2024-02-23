import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import React from "react"
import { ChangePassword, DifficultyPage, LoginPage, MatchingPage, QuestionPage, SignupPage } from "./view";
import PrivateRoutes from "./utils/PrivateRoutes";
import NotLoggedInRoutes from "./utils/NotLoggedInRoutes";
import InMatchRoutes from "./utils/InMatchRoutes";
import NotInMatchRoutes from "./utils/NotInMatchRoutes";

function App() {
    return (
        <div className="App">
            <Router>
                <Routes>
                    <Route element={<NotLoggedInRoutes/>}>
                        <Route exact path="/" element={<Navigate replace to="/login" />} />
                        <Route path="/login" element={<LoginPage/>}/>
                        <Route path="/signup" element={<SignupPage />}/>
                    </Route>
                    <Route element={<PrivateRoutes/>}>
                        <Route element={<NotInMatchRoutes/>}>
                            <Route path="/difficulty" element={<DifficultyPage/>}/>
                            <Route path="/change-password" element={<ChangePassword/>}/>
                            <Route path="/matching" element={<MatchingPage/>}/>
                        </Route>
                        <Route element={<InMatchRoutes/>}>
                            <Route path="/question" element={<QuestionPage/>}/>
                        </Route>
                    </Route>
                </Routes>
            </Router>
        </div>
    );
}

export default App;
