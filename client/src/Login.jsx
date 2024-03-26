import axios from 'axios';
import React, { useState } from 'react'

export default function Login() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const pressSubmit = async () => {
        console.log(email, password);
        try {
            const res = await axios.post("http://localhost:5000/api/user/login", {
                email: email,
                password: password
            })
            console.log(res.data);
            if (res.status === 202) {
                localStorage.setItem('email', res.data.email);
                localStorage.setItem('_id', res.data._id);
                window.location.href = '/allproducts';
            }
        } catch (error) {
            console.log(error);
        }
    }

    return (
        <div className='w-1/2 h-full mx-auto '>
            <div className='font-bold text-3xl mt-20'>Login</div>
            <hr />
            <div className='mt-10'>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Email</div>
                    <input type='text' className='border-2' onChange={
                        (e) => {
                            // console.log(e.target.value);
                            setEmail(e.target.value);
                        }
                    } />
                </div>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Password</div>
                    <input type='password' className='border-2' onChange={
                        (e) => {
                            // console.log(e.target.value);
                            setPassword(e.target.value);
                        }
                    } />
                </div>
                <button className='bg-green-500 text-white px-2 py-1 rounded m-2' type='button' onClick={pressSubmit}>Login</button>
            </div>
        </div>
    )
}