import React, {useContext, useEffect, useState} from 'react';
import {useHttp} from "../hooks/http.hook";
import {AuthContext} from "../context/AuthContext";

export const AuthPage = () => {
    const auth = useContext(AuthContext)
    const {loading, error, request, clearError} = useHttp()
    const [form, setForm] = useState({
        username: '', password: ''
    })

    const changeHandler = event => {
        setForm({ ...form, [event.target.name]: event.target.value })
    }

    const loginHandler = async (event) => {
        event.preventDefault()
        clearError()
        try {
            const data = await request('http://localhost:8080/api/login', 'POST', {...form})
            auth.login(data.token)
        } catch (e) {
        }
    }

    return (
        <>
            <h1 className="text-center mt-5">Увійти в систему Друг</h1>
            <form>
                <div className="form-group">
                    <label htmlFor="username">Ім'я користувача</label>
                    <input onChange={changeHandler} type="text" className={"form-control" + (error ? " is-invalid" : "")} id="username" name="username"/>
                </div>
                <div className="form-group">
                    <label htmlFor="password">Пароль</label>
                    <input onChange={changeHandler} type="password" className={"form-control" + (error ? " is-invalid" : "")} id="password" name="password"/>
                </div>
                <button type="submit" className="btn btn-primary" onClick={loginHandler} disabled={loading}>Увійти</button>
            </form>
        </>
    )
}