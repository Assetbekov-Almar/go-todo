'use client';

import axios from 'axios';
import { useRouter } from 'next/navigation';
import qs from 'qs';
import { toast } from 'react-toastify';

const Register = () => {
    const router = useRouter();
    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const formData = new FormData(e.target as HTMLFormElement);
        const username = formData.get('username');
        const password = formData.get('password');
        axios
            .post(
                'http://localhost:8082/register',
                qs.stringify({
                    username,
                    password,
                }),
                { headers: { 'Content-Type': 'application/x-www-form-urlencoded' } }
            )
            .then(() => {
                router.push('/login');
            })
            .catch((error) => {
                toast.error(error.response.data.error);
                return;
            });
    };

    return (
        <div className="flex max-h-screen flex-col items-center justify-between p-24">
            <form className="flex flex-col items center gap-6" onSubmit={onSubmit}>
                <div className="flex flex-col gap-2">
                    <label htmlFor="username">Username</label>
                    <input type="username" name="username" />
                </div>
                <div className="flex flex-col gap-2">
                    <label htmlFor="password">Password</label>
                    <input type="password" name="password" />
                </div>
                <button type="submit" className="bg-lime-600 p-2 rounded">
                    Register
                </button>
            </form>
        </div>
    );
};

export default Register;
