import React from 'react'
import {Switch, Route, Redirect} from "react-router-dom";
import {QuestionsPage} from "./pages/QuestionsPage";
import {QuestionDetailPage} from "./pages/QuestionDetailPage";
import {QAPage} from "./pages/QAPage";
import {QADetailPage} from "./pages/QADetailPage";
import {AuthPage} from "./pages/AuthPage";
import {QACreatePage} from "./pages/QACreatePage";
import {SettingsPage} from "./pages/SettingsPage";

export const useRoutes = isAuthenticated => {
    if (isAuthenticated) {
        return (
            <Switch>
                <Route path="/questions" exact>
                    <QuestionsPage />
                </Route>
                <Route path="/question/:id">
                    <QuestionDetailPage />
                </Route>
                <Route path="/qa" exact>
                    <QAPage />
                </Route>
                <Route path="/qa/create" exact>
                    <QACreatePage />
                </Route>
                <Route path="/qa/:id">
                    <QADetailPage />
                </Route>
                <Route path="/settings" exact>
                    <SettingsPage />
                </Route>
                <Redirect to="/questions" />
            </Switch>
        )
    }

    return (
        <Switch>
            <Route path="/" exact>
                <AuthPage />
            </Route>
            <Redirect to="/" />
        </Switch>
    )
}