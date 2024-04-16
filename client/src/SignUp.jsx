import axios from 'axios';
import React, {useState} from 'react'
import { useNavigation } from 'react-router-dom';

export default function SignUp() {

    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');

    const handleSubmit = async () => {
        if (password !== confirmPassword) {
            alert('Passwords do not match');
            return;
        }
        if (name === '' || email === '' || password === '' || confirmPassword === '') {
            alert('All fields are required');
            return;
        }
        try {
            const res = await axios.post('http://localhost:5000/api/user/create', {
                username: name,
                email: email,
                password: password
            });
            console.log(res);
            if (res.status === 201) {
                alert('User created successfully');
                // redirect to login page
                window.location.href = '/login';         
            }
        } catch (error) {
            console.log(error);
            alert('Error in creating user');
        }
    }

    return (
        <div className='w-1/2 h-full mx-auto '>
            <div className='font-bold text-3xl mt-20'>SignUp</div>
            <hr />
            <div className='mt-10'>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Name</div>
                    <input type='text' className='border-2' 
                        onChange={(e) => {
                            setName(e.target.value);
                        }}
                    />
                </div>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Email</div>
                    <input type='text' className='border-2' 
                        onChange={(e) => {
                            setEmail(e.target.value);
                        }}
                    />
                </div>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Password</div>
                    <input type='password' className='border-2' 
                        onChange={(e) => {
                            setPassword(e.target.value);
                        }}
                    />
                </div>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Confirm Password</div>
                    <input type='password' className='border-2' 
                        onChange={(e) => {
                            setConfirmPassword(e.target.value);
                        }}
                    />
                </div>
                <button className='bg-green-500 text-white px-2 py-1 rounded m-2' type='button' onClick={handleSubmit}>Sign Up</button>
            </div>
        </div>
    )
}
