import React from 'react'

export default function SignUp() {
    return (
        <div className='w-1/2 h-full mx-auto '>
            <div className='font-bold text-3xl mt-20'>SignUp</div>
            <hr />
            <div className='mt-10'>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Name</div>
                    <input type='text' className='border-2' />
                </div>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Email</div>
                    <input type='text' className='border-2' />
                </div>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Password</div>
                    <input type='password' className='border-2' />
                </div>
                <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
                    <div>Confirm Password</div>
                    <input type='password' className='border-2' />
                </div>
                <button className='bg-green-500 text-white px-2 py-1 rounded m-2' type='button'>Sign Up</button>
            </div>
        </div>
    )
}
