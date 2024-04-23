import React from 'react'
import axios from 'axios'

export default function Profile() {
  const [products, setProducts] = React.useState([])
  const [loading, setLoading] = React.useState(true)

  React.useEffect(() => {
    const f = async () => {
      const res = await axios.post("http://localhost/api/user/previouslyBought", {
        email: localStorage.getItem('email')
      })
      console.log(res.data)
      setProducts(res.data)
      setLoading(false)
    }
    f()
  }, [])


  return (
    <div className='w-1/2 h-full mx-auto '>
      <div className='font-bold text-3xl mt-20'>{localStorage.getItem("email")}</div>
      <hr />
      <div className='font-bold text-xl mt-5'>Products purchase history</div>
      <div className='mt-10'>
        {loading ? <div>LOADING...</div> : products.map((product) => (
          <div key={product.id} className='flex justify-between border-b-2 p-3 m-2 bg-gray-300 '>
            <div>{product.name}</div>
            <div>Price: {product.price}</div>
          </div>
        ))}
      </div>
    </div>
  )
}
