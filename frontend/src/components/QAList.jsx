import React, {useEffect, useState} from 'react'
import {useHistory} from "react-router-dom";

export const QAList = ({qas, topics}) => {
    const history = useHistory()
    const [qasToShow, setQAsToShow] = useState([])
    const [topic, setTopic] = useState(0)

    useEffect(() => {
        setQAsToShow([...qas])
    }, [qas])

    useEffect(() => {
        if (topic === 0) {
            setQAsToShow([...qas])
        } else {
            setQAsToShow(qas.filter(qa => qa.topic_id === topic))
        }
    }, [topic, topics, qas])

    const handleClick = (id) => {
        history.push(`/qa/${id}`)
    }

    const listItems = qasToShow.map((qa) => {
        return (
            <tr key={qa.id} onClick={() => handleClick(qa.id)}>
                <td>{qa.id}</td>
                <td>{qa.topic_name}</td>
                <td>{qa.question}</td>
            </tr>
        )
    })

    const topicItems = topics.map((topic) => {
        return (
            <option key={topic.id} value={topic.id}>
                {topic.name}
            </option>
        )
    })

    const handleTopic = (event) => {
        setTopic(parseInt(event.target.value))
    }

    const handleAdd = (event) => {
        event.preventDefault()
        history.push('/qa/create')
    }

    return (
        <>
            <div style={{marginBottom: '10px', marginTop: '10px'}}>
                <button style={{marginBottom: '10px', float: 'right'}} className="btn btn-success" onClick={handleAdd}>Додати</button>
                <select className="form-control" value={topic} onChange={handleTopic}>
                    <option value="0">Всі теми</option>
                    {topicItems}
                </select>
            </div>
            {
                qasToShow.length === 0 ? (
                    <h1>Питання відсутні</h1>
                ) : (
                    <table className="table table-light table-hover">
                        <thead className="thead-dark">
                        <tr>
                            <th style={{width: '10%'}}>№</th>
                            <th style={{width: '30%'}}>Тема</th>
                            <th style={{width: '60%'}}>Питання</th>
                        </tr>
                        </thead>
                        <tbody>
                            {listItems}
                        </tbody>
                    </table>
                )

            }
        </>
    )
}