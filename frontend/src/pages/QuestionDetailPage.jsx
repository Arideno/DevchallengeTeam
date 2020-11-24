import React, {useCallback, useContext, useEffect, useState} from 'react';
import {useParams} from 'react-router-dom'
import {useHttp} from "../hooks/http.hook";
import {AuthContext} from "../context/AuthContext";
import '../css/QuestionDetailPage.css'

export const QuestionDetailPage = () => {
    const [question, setQuestion] = useState({})
    const {token, logout} = useContext(AuthContext)
    const {loading, request} = useHttp()
    const questionId = useParams().id

    const getQuestion = useCallback(async () => {
        try {
            const data = await request(`http://localhost:8080/api/auth/questions/${questionId}`, 'GET', null, {
                Authorization: `Bearer ${token}`
            })
            setQuestion(data)
        } catch (e) {
        }
    }, [token, request, questionId])

    const changeStatus = async (event) => {
        const status = parseInt(event.target.value)
        setQuestion({...question, status: status})
        await request(`http://localhost:8080/api/auth/changeStatus`, 'PATCH', {
            "status": status,
            "question_id": parseInt(questionId)
        }, {
            Authorization: `Bearer ${token}`
        })
    }

    useEffect(() => {
        getQuestion()
    }, [getQuestion])

    if (loading) {
        return (
            <h1>Loading</h1>
        )
    }

    return (
        <>
        { !loading &&
            <>
                <div className="jumbotron">
                    <h1 className="display-5">Питання №{question.id}</h1>
                    <p className="lead">{question.question}</p>
                    <span style={{marginRight: '5px'}}>Статус</span>
                    <select className="form-control-sm" value={question.status} onChange={changeStatus}>
                        <option value="0">Не вирішено</option>
                        <option value="1">В процесі</option>
                        <option value="2">Вирішено</option>
                    </select>
                </div>
                <div>
                    <ul className="chat">
                        <li className="message mine">Message1</li>
                        <li className="message">Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1Message1</li>
                        <li className="message">Message1</li>
                        <li className="message mine">Message1</li>
                        <li className="message">Message1</li>
                        <li className="message">Message1</li>
                    </ul>
                </div>
            </>
        }
        </>
    )
}