import React, {useCallback, useContext, useEffect, useState} from 'react';
import {useHttp} from "../hooks/http.hook";
import {QuestionList} from "../components/QuestionList";
import {AuthContext} from "../context/AuthContext";
import {useHistory} from "react-router-dom";

export const QuestionsPage = () => {
    const history = useHistory()
    const {token, logout} = useContext(AuthContext)
    const {loading, request} = useHttp()

    const [questions, setQuestions] = useState([])

    const getQuestions = useCallback(async () => {
        try {
            const data = await request('http://localhost:8080/api/auth/questions', 'GET', null, {
                Authorization: `Bearer ${token}`
            })
            setQuestions(data)
        } catch (e) {
        }
    }, [token, logout, request])

    useEffect(() => {
        getQuestions()
    }, [getQuestions])

    return (
        <div>
            <QuestionList questions={questions}/>
        </div>
    )
}