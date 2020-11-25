import React, {useContext, useState} from 'react'
import {AuthContext} from "../context/AuthContext";
import {useHttp} from "../hooks/http.hook";
import {useHistory} from 'react-router-dom'

export const SettingsPage = () => {
    const history = useHistory()
    const [form, setForm] = useState({
        currentPassword: '',
        newPassword: '',
        repeatPassword: '',
        match: false
    })
    const [topics, setTopics] = useState([])
    const {token, logout} = useContext(AuthContext)
    const {loading, request} = useHttp()

    const handleChange = (event) => {
        let match = true;
        if (event.target.name === 'newPassword') {
            match = (event.target.value === form.repeatPassword) && (event.target.value !== '')
        } else if (event.target.name === 'repeatPassword') {
            match = event.target.value === form.newPassword && (event.target.value !== '')
        } else {
            match = (form.repeatPassword === form.newPassword) && (form.newPassword !== '')
        }

        setForm({
            ...form,
            [event.target.name]: event.target.value,
            match: match
        })
    }

    const handleClick = async (event) => {
        event.preventDefault()

        try {
            const data = await request('http://localhost:8080/api/auth/change/password', 'PATCH', {
                current_password: form.currentPassword,
                new_password: form.newPassword,
            }, {
                Authorization: `Bearer ${token}`
            })
            if (data.message === "ok") {
                history.push('/questions')
            }
        } catch (e) {
            const parsedError = JSON.parse(e.message)
            if (parsedError.code === 401) {
                logout()
            }
        }
    }

    return (
        <div>
            <h1 className="text-center">Налаштування</h1>
            <form>
                <div className="form-group">
                    <label htmlFor="currentPassword">Поточний пароль</label>
                    <input className="form-control" type="password" id="currentPassword" name="currentPassword" value={form.currentPassword} onChange={handleChange}/>
                </div>
                <div className="form-group">
                    <label htmlFor="newPassword">Новий пароль</label>
                    <input className={`form-control ${form.match ? "" : "is-invalid"}`} type="password" id="newPassword" name="newPassword" value={form.newPassword} onChange={handleChange}/>
                </div>
                <div className="form-group">
                    <label htmlFor="repeatPassword">Повторіть новий пароль</label>
                    <input className={`form-control ${form.match ? "" : "is-invalid"}`} type="password" id="repeatPassword" name="repeatPassword" value={form.repeatPassword} onChange={handleChange}/>
                </div>
                <button type="submit" className="btn btn-primary" onClick={handleClick} disabled={!form.match}>Зберегти</button>
            </form>
        </div>
    )
}