import { useState } from "react";
import { Login, Logout } from "./utils";

export default function Home() {
  const [data, setData] = useState();

  const[login,setLogin] = useState({email:"",password:""})
  const [validity, setValidity] = useState({ email: true, password: true });


  // Event handler for input change
  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setLogin({ ...login, [name]: value });
    // validation of input data
    setValidity({ ...validity, [name]: e.target.checkValidity() });
  };

  // Event handler for form submission
  const handleSubmit = (e) => {
    e.preventDefault();
    // Check if all inputs are valid
    const isFormValid = Object.values(validity).every((isValid) => isValid);
    if (isFormValid) {
        Login(login.email,login.password,setData)
          // Clear the input fields after submission
          setLogin({ email: '', password: '' });
          // Reset input validity
          setValidity({ email: true, password: true });
        } else {
          console.log('Form is invalid. Please correct the errors.');
        }
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
        className={validity.email ? '' : 'invalid'}
        placeholder="Email"
        name="email"
        value={login.email}
        onChange={handleInputChange}
        required
      />
      {!validity.email && (
        <p className="invalid-message">Please enter a valid email address.</p>
      )}
      <input
        type="password"
        className={validity.password ? '' : 'invalid'}
        placeholder="Password"
        name="password"
        value={login.password}
        onChange={handleInputChange}
        required
      />
      {!validity.password && (
        <p className="invalid-message">
          Password must be at least 8 characters long.
        </p>
      )}
      <button type="submit">Submit</button>
    </form>
    }
      <button className="px-4 py-2 bg-red-400 text-white" onClick={()=>Logout(setData)}>Logout</button>
    </main>
  );
}