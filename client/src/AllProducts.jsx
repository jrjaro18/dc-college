import React, { useState, useEffect } from 'react'
import axios from 'axios';
// import Cart from './Cart';
export default function AllProducts() {
    // const products = [
    //     { id: 1, name: 'Product 1', price: 100 },
    //     { id: 2, name: 'Product 2', price: 200 },
    //     { id: 3, name: 'Product 3', price: 300 },
    //     { id: 4, name: 'Product 4', price: 400 },
    //     { id: 5, name: 'Product 5', price: 500 },
    //     { id: 6, name: 'Product 6', price: 600 },
    //     { id: 7, name: 'Product 7', price: 700 },
    //     { id: 8, name: 'Product 8', price: 800 },
    //     { id: 9, name: 'Product 9', price: 900 },
    //     { id: 10, name: 'Product 10', price: 1000 },
    // ];

    const [products, setProducts] = useState([]);
    useEffect(() => {
        const f = async () => {
            try {
                const res = await axios.get('http://localhost/api/item', {})
                console.log(res.data);
                setProducts(res.data);
            } catch (error) {
                console.log(error);
            }
        }
        f()
    }, []);

    const handleAddToCart = async (id) => {
        try {
            const res = await axios.post('http://localhost/api/user/add', {
                email: localStorage.getItem('email'),
                itemID: id
            })
            console.log(res.data);
            if (res.data.status === 202) {
                alert('Item added to cart');
            }

        } catch (error) {
            alert("error in adding item to cart");
        }
    }

    return (
        <div className='w-1/2 h-full mx-auto'>
            <div className='font-bold text-3xl mt-20'>All Products</div>
            <hr />
            <div className='mt-10'>
                {products.map((product) => (
                    <div key={product._id} className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '
                        onClick={() => {
                            handleAddToCart(product._id);
                        }}
                    >
                        <div>{product.name}</div>
                        <div>Price: {product.price}</div>
                        <button className='bg-green-500 text-white px-2 py-1 rounded' type='button'>Add to Cart</button>
                    </div>
                ))}
            </div>
        </div>
    )
}
