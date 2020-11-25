import React from 'react'
import {useHistory} from "react-router-dom";

export const QuestionList = ({questions}) => {
    const history = useHistory()

    const handleClick = (id) => {
        history.push(`/question/${id}`)
    }

    const listItems = questions.map((question) => {
        let status;
        if (question.status === 0) {
            status = "Не вирішено 🔴"
        } else if (question.status === 1) {
            status = "В процесі 🟡"
        } else {
            status = "Вирішено 🟢"
        }
        return (
            <tr key={question.id} onClick={() => handleClick(question.id)}>
                <td>{question.id}</td>
                <td>{question.question}</td>
                <td>{status}</td>
            </tr>
        )
    })

    if (questions.length === 0) {
        return <h1>Питання відсутні</h1>
    }

    return (
        <table className="table table-light table-hover">
            <thead className="thead-dark">
                <tr>
                    <th style={{width: '10%'}}>№</th>
                    <th style={{width: '70%'}}>Питання</th>
                    <th style={{width: '20%'}}>Статус</th>
                </tr>
            </thead>
            <tbody>
                {listItems}
            </tbody>
        </table>
    )
}