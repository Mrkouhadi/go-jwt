import { useState } from "react";
import { Login, Logout } from "./utils";

export default function Home() {
  const [data, setData] = useState();
  const[login,setLogin] = useState({email:"",password:""})


  // Event handler for input change
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setLogin({ ...login, [name]: value });
  };

  // Event handler for form submission
  const handleSubmit = (e) => {
    e.preventDefault();
    Login(login.email,login.password,setData)
    // Clear the input fields after submission
    setLogin({ email: '', password: '' });
  };
  return (
    <main className="flex h-screen flex-col items-center gap-12 p-12 bg-pink-400 text-black">
      {data ? (
        <>
          <p>email: {data.data.email}</p>
          <p>username: {data.data.username}</p>
          <p>role: {data.data.role}</p>
        </>
      ):
      <form onSubmit={handleSubmit}>
          <input
            type="email"
            className=""
            placeholder="Email"
            name="email"
            value={login.email}
            onChange={handleInputChange}
          />
          <input
            type="password"
            className=""
            placeholder="Password"
            name="password"
            value={login.password}
            onChange={handleInputChange}
          />
      <button type="submit">Login</button>
    </form>
    }
      <button className="px-4 py-2 bg-red-400 text-white" onClick={()=>Logout(setData)}>Logout</button>
    </main>
  );
}