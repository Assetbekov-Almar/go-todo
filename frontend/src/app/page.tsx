'use client';

import axios from 'axios';
import Image from 'next/image';
import { useEffect, useState } from 'react';

const fetchTodos = async () => {
    const { data } = await axios.get('http://localhost:8080/todo/all');
    return data;
};

export default function Home() {
    const [todos, setTodos] = useState([]);
    const [editedTodoId, setEditedTodoId] = useState(null);

    useEffect(() => {
        fetchTodos().then((data) => setTodos(data));
    }, []);

    return (
        <main className="flex min-h-screen flex-col items-center justify-between p-24">
            <div className="flex flex-col items-center space-y-12 justify-center">
                {todos.map((todo) => (
                    <form
                        key={todo.id}
                        className="flex space-x-8 justify-center items-center text-black"
                        onSubmit={async (e) => {
                            e.preventDefault();
                            if (editedTodoId === todo.id) {
                                const formData = new FormData(e.target as HTMLFormElement);
                                const title = formData.get('title');
                                const description = formData.get('description');
                                const deadline = formData.get('deadline');
                                if (title === todo.title && description === todo.description && deadline === todo.deadline) {
                                    setEditedTodoId(null);
                                    return;
                                }

                                const result = await axios.patch(`http://localhost:8080/todo/update/${todo.id}`, {
                                    description,
                                    deadline,
                                });

                                const data = await fetchTodos();
                                setTodos(data);

                                setEditedTodoId(null);
                                return;
                            }
                        }}
                    >
                        <fieldset disabled={todo.id !== editedTodoId} className="flex space-x-8 justify-center items-center text-black ">
                            <input name="title" defaultValue={todo.title} />
                            <input name="description" defaultValue={todo.description} />
                            <input name="deadline" defaultValue={todo.deadline} />
                        </fieldset>
                        {editedTodoId === todo.id ? (
                            <button type="submit">
                                <Image src={'/done.svg'} alt="edit" width="28" height="28" />
                            </button>
                        ) : (
                            <button
                                onClick={(e) => {
                                    e.preventDefault();
                                    setEditedTodoId((prevState) => (prevState === todo.id ? null : todo.id));
                                }}
                            >
                                <Image src={'/pen.svg'} alt="edit" width="28" height="28" />
                            </button>
                        )}
                    </form>
                ))}
            </div>
        </main>
    );
}
