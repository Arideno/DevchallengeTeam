import React, {useCallback, useContext, useEffect, useState} from 'react';
import {QAList} from "../components/QAList";
import {useHistory} from "react-router-dom";
import {AuthContext} from "../context/AuthContext";
import {useHttp} from "../hooks/http.hook";

export const QAPage = () => {
    const history = useHistory()
    const {token, logout} = useContext(AuthContext)
    const {loading, request} = useHttp()

    const [qas, setQAs] = useState([])
    const [topics, setTopics] = useState([])

    const getQAs = useCallback(async () => {
        try {
            const data = await request('http://localhost:8080/api/auth/qa', 'GET', null, {
                Authorization: `Bearer ${token}`
            })
            setQAs(data)
        } catch (e) {
            const parsedError = JSON.parse(e.message)
            if (parsedError.code === 401) {
                logout()
            }
        }
    }, [token, logout, request])

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
        getQAs()
        getTopics()
    }, [getQAs])

    return (
        <div>
            <QAList qas={qas} topics={topics} />
        </div>
    )
}