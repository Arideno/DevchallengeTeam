import React, {useContext} from 'react'
import {NavLink, useHistory} from 'react-router-dom'
import {AuthContext} from "../context/AuthContext";

export const Navbar = () => {
    const history = useHistory()
    const auth = useContext(AuthContext)

    const logoutHandler = event => {
        event.preventDefault()
        auth.logout()
        history.push('/')
    }

    return (
        <nav className="navbar navbar-expand-sm navbar-light bg-light">
            <a className="navbar-brand" href="#">Друг</a>
            <ul className="navbar-nav mr-auto">
                <li className="nav-item">
                    <NavLink to="/questions" className="nav-link">Питання</NavLink>
                </li>
                <li className="nav-item">
                    <NavLink to="/qa" className="nav-link">База знань</NavLink>
                </li>
                <li className="nav-item">
                    <NavLink to="/settings" className="nav-link">Налаштування</NavLink>
                </li>
                <li className="nav-item">
                    <a href="#" onClick={logoutHandler} className="nav-link">Вийти</a>
                </li>
            </ul>
        </nav>
    )
}