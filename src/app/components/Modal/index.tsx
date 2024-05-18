import axios from 'axios';
import Image from 'next/image';
import qs from 'qs';
import { toast } from 'react-toastify';

interface ModalProps {
    close: () => void;
    todoCreated: () => void;
}

const Modal = ({ close, todoCreated }: ModalProps) => {
    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        const formData = new FormData(e.target as HTMLFormElement);
        const object = Object.fromEntries(formData.entries());
        const { title, description, deadline } = object;

        if (!title || !description || !deadline) {
            toast.error('Please fill all fields');
            return;
        }

        try {
            const { status } = await axios.post(`https://todo-go.fly.dev/todo/add`, qs.stringify(object), {
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            });

            if (status !== 200) {
                toast.error('Error creating todo');
                return;
            }

            toast.success('Todo created');
            todoCreated();
            close();
        } catch (error) {
            console.error(error);
            toast.error('Error creating todo');
            return;
        }
    };

    return (
        <>
            <div className="absolute z-10 w-full h-full bg-black opacity-70 cursor-pointer" onClick={close} />
            <div className="absolute z-20 top-1/3 left-1/2 transform -translate-x-1/2 -translate-y-1/3 flex flex-col items-center bg-slate-400 w-[400px] h-[600px] rounded-lg">
                <Image
                    src="/close.svg"
                    alt="close"
                    width="28"
                    height="28"
                    onClick={close}
                    className="cursor-pointer align-right right-0 absolute h-10 w-10 top-2 right-2"
                />
                <form className="flex flex-col justify-center items-center h-full" onSubmit={handleSubmit}>
                    <fieldset className="flex flex-col space-y-6 justify-center items-center text-black">
                        <input name="title" placeholder="Title" className="w-full" />
                        <textarea name="description" placeholder="Description" rows={8} cols={38} className="rounded-lg" />
                        <input
                            name="deadline"
                            type="date"
                            placeholder="Deadline"
                            className="w-full"
                            defaultValue={new Date().toISOString().substr(0, 10)}
                        />
                    </fieldset>
                    <button
                        type="submit"
                        className="bg-green-700 text-white absolute rounded-b-lg h-10 w-10 bottom-0 left-1/2 transform -translate-x-1/2 w-full text-xl hover:bg-green-600 transition ease-in-out delay-100"
                    >
                        Create
                    </button>
                </form>
            </div>
        </>
    );
};

export default Modal;
