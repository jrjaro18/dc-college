import axios from 'axios'
import React from 'react'

export default function Cart() {
    const [products, setProducts] = React.useState([])
    const [loading, setLoading] = React.useState(true)

    React.useEffect(() => {
        const f = async () => {
            const res = await axios.post("http://localhost:5000/api/user/cart", {
                email: localStorage.getItem('email')
            })
            console.log(res.data)
            setProducts(res.data)
            setLoading(false)
        }
        f()
    }, [])

    const onSubmitCart = async () => {
        try{
            const res = await axios.post("http://localhost:5000/api/user/buy", {
                email: localStorage.getItem('email'),
            })
            console.log(res.data)
            if(res.status === 202){
                alert('Successfully Bought')
                window.location.reload()
            }
        } catch(err){
            console.log(err)
        }
    }


    return (
        <div className='w-1/2 h-full mx-auto '>
            <div className='font-bold text-3xl mt-20'>Cart List</div>
            <hr />
            <div className='mt-10'>
                {
                    loading ? <div>LOADING...</div> : products.map((product) => (
                        <div key={product._id} className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '

                        >
                            <div>{product.name}</div>
                            <div>Price: {product.price}</div>
                        </div>
                    ))
                }
            </div>
            <div className='flex justify-between border-b-2 p-3 m-2 bg-gray-600 text-white'>
                <div>Total Price: </div>
                <div>{products.reduce((acc, product) => acc + product.price, 0)}</div>
            </div>

            <button className='bg-green-500 text-white px-2 py-1 rounded m-2' type='button'
                onClick={onSubmitCart}
            >Buy Now</button>
        </div>
    )
}
