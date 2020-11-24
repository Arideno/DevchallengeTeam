import React, {useCallback, useContext, useEffect, useRef, useState} from 'react';
import {useParams} from 'react-router-dom'
import {useHttp} from "../hooks/http.hook";
import {AuthContext} from "../context/AuthContext";
import '../css/QuestionDetailPage.css'

export const QuestionDetailPage = () => {
    const [question, setQuestion] = useState({})
    const [messages, setMessages] = useState([])
    const [text, setText] = useState("")
    const {token, logout} = useContext(AuthContext)
    const {loading, request} = useHttp()
    const questionId = useParams().id
    const socket = useRef(null)
    const chatRef = useRef(null)

    useEffect(() => {
        socket.current = new WebSocket('ws://localhost:8080/ws')
        socket.current.onopen = () => {
            setTimeout(() => {
                socket.current.send(JSON.stringify({
                type: 'GET_MESSAGES',
                data: {
                    questionId: parseInt(questionId)
                },
                token: token
            }))
            }, 1000)
        }
    }, [])

    const scrollToBottom = () => {
        chatRef.current.scrollTop = chatRef.current.scrollHeight
    }

    useEffect(() => {
        socket.current.onmessage = (message) => {
            if (message.data === "0") {
                console.log("error")
            } else {
                const data = JSON.parse(message.data)
                if (data.type === "message") {
                    setMessages([...messages, data.data])
                    setText("")
                } else if (data.type === "messages") {
                    setMessages(data.data)
                }
            }
        }
        scrollToBottom()
    }, [messages])

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

    const handleChange = (event) => {
        setText(event.target.value)
    }

    const handleSend = () => {
        if (text && text.trim()) {
            socket.current.send(JSON.stringify({
                type: 'SEND_MESSAGE',
                data: {
                    message: text,
                    chatId: question.chat_id,
                    questionId: question.id,
                },
                token: token
            }))
        }
    }

    useEffect(() => {
        getQuestion()
    }, [getQuestion])

    if (loading) {
        return (
            <h1>Loading</h1>
        )
    }

    const messageItems = messages.map((message) => {
        return (
            <li key={message.id} className={"message" + (message.from_operator ? " mine" : "")}>{message.message}</li>
        )
    })

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
                    <ul className="chat" ref={chatRef}>
                        {messageItems}
                    </ul>
                </div>
                <div style={{display: 'flex', alignItems: 'center'}}>
                    <textarea className="input" value={text} onChange={handleChange}/>
                    <button className="send" onClick={handleSend}>
                        <svg width="50px" height="50px" viewBox="0 0 16 16" className="bi bi-arrow-right-circle-fill"
                             fill="#8B8B8B" xmlns="http://www.w3.org/2000/svg">
                            <path fillRule="evenodd" d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zm-11.5.5a.5.5 0 0 1 0-1h5.793L8.146 5.354a.5.5 0 1 1 .708-.708l3 3a.5.5 0 0 1 0 .708l-3 3a.5.5 0 0 1-.708-.708L10.293 8.5H4.5z"/>
                        </svg>
                    </button>
                </div>
            </>
        }
        </>
    )
}