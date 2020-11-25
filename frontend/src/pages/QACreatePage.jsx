import React, {useCallback, useContext, useEffect, useState} from 'react'
import {useHistory} from "react-router-dom";
import {AuthContext} from "../context/AuthContext";
import {useHttp} from "../hooks/http.hook";

export const QACreatePage = () => {
    const history = useHistory()
    const [form, setForm] = useState({
        topic: 0,
        question: '',
        answer: ''
    })
    const [topics, setTopics] = useState([])
    const {token, logout} = useContext(AuthContext)
    const {loading, request} = useHttp()

    const getTopics = useCallback(async () => {
        try {
            const data = await request('http://localhost:8080/api/auth/topics', 'GET', null, {
                Authorization: `Bearer ${token}`
            })
            setTopics(data)
        } catch (e) {
            const parsedError = JSON.parse(e.message)
            if (parsedError.code === 401) {
                logout()
            }
        }
    }, [])

    useEffect(() => {
        getTopics()
    }, [])

    const topicItems = topics.map((topic) => {
        return (
            <option key={topic.id} value={topic.id}>
                {topic.name}
            </option>
        )
    })

    const createQA = async () => {
        try {
            const data = await request(`http://localhost:8080/api/auth/qa`, 'POST', {
                topic_id: parseInt(form.topic),
                question: form.question,
                answer: form.answer,
            }, {
                Authorization: `Bearer ${token}`
            })
            if (data.message === "ok") {
                history.push('/qa')
            }
        } catch (e) {
            const parsedError = JSON.parse(e.message)
            if (parsedError.code === 401) {
                logout()
            }
        }
    }

    const handleChange = (event) => {
        setForm({...form, [event.target.name]: event.target.value})
    }

    const handleClick = (event) => {
        event.preventDefault()
        createQA()
    }

    return (
        <>
            <form>
                <div className="form-group">
                    <label htmlFor="topic">Тема</label>
                    <select className="form-control" id="topic" name="topic" value={form.topic} onChange={handleChange}>
                        {topicItems}
                    </select>
                </div>
                <div className="form-group">
                    <label htmlFor="question">Питання</label>
                    <input type="text" className="form-control" id="question" name="question" value={form.question} onChange={handleChange}/>
                </div>
                <div className="form-group">
                    <label htmlFor="answer">Відповідь</label>
                    <textarea className="form-control" id="answer" name="answer" value={form.answer} onChange={handleChange}/>
                </div>
                <button type="submit" className="btn btn-primary" onClick={handleClick}>Створити</button>
            </form>
        </>
    )
}