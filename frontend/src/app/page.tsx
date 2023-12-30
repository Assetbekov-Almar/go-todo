'use client';

import axios from 'axios';
import { format } from 'date-fns';
import Image from 'next/image';
import { useEffect, useState } from 'react';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

interface Todo {
    id: number;
    title: string;
    description: string;
    deadline: string;
}

const fetchTodos = async (): Promise<Todo[]> => {
    try {
        const { data } = await axios.get('http://localhost:8080/todo/all');
        return data;
    } catch (error) {
        console.error(error);
        return [];
    }
};

export default function Home() {
    const [todos, setTodos] = useState<Todo[]>([]);
    const [editedTodoId, setEditedTodoId] = useState<number | null>(null);

    useEffect(() => {
        fetchTodos().then((data) => setTodos(data));
    }, []);

    const submitHandler = async (e: React.FormEvent<HTMLFormElement>, todo: Todo) => {
        e.preventDefault();
        if (editedTodoId === todo.id) {
            const formData = new FormData(e.target as HTMLFormElement);
            const title = formData.get('title');
            const description = formData.get('description');
            const deadline = formData.get('deadline');
            if (title === todo.title && description === todo.description && deadline === format(todo.deadline, 'MM/dd/yyyy')) {
                setEditedTodoId(null);
                return;
            }

            try {
                const { status } = await axios.patch(`http://localhost:8080/todo/update/${todo.id}`, {
                    description,
                    deadline,
                });

                if (status !== 200) {
                    toast.error('Error updating todo');
                    return;
                }

                toast.success('Todo updated');

                const data = await fetchTodos();
                setTodos(data);
            } catch (error) {
                console.error(error);
                toast.error('Error updating todo');
                return;
            }

            setEditedTodoId(null);
            return;
        }
    };

    return (
        <main>
            <ToastContainer />
            <div className="flex min-h-screen flex-col items-center justify-between p-24">
                <div className="flex flex-col items-center space-y-12 justify-center">
                    {todos.map((todo) => (
                        <form
                            key={todo.id}
                            className="flex space-x-8 justify-center items-center text-black"
                            onSubmit={(e) => submitHandler(e, todo)}
                        >
                            <fieldset
                                disabled={todo.id !== editedTodoId}
                                className="flex space-x-8 justify-center items-center text-black "
                            >
                                <input name="title" defaultValue={todo.title} />
                                <input name="description" defaultValue={todo.description} />
                                <input name="deadline" defaultValue={format(todo.deadline, 'MM/dd/yyyy')} />
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
                            <button
                                onClick={async (e) => {
                                    e.preventDefault();
                                    const { status } = await axios.delete(`http://localhost:8080/todo/delete/${todo.id}`);
                                    if (status !== 200) {
                                        toast.error('Error deleting todo');
                                        return;
                                    }

                                    toast.success('Todo deleted');

                                    const data = await fetchTodos();
                                    setTodos(data);
                                }}
                            >
                                <Image src={'/remove.svg'} alt="edit" width="28" height="28" />
                            </button>
                        </form>
                    ))}
                </div>
            </div>
        </main>
    );
}
