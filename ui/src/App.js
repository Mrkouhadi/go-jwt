import { useState } from "react";
import Profile from "./components/Profile";
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
        <Profile toLogout={()=>Logout(setData)} email={data.data.email} username={data.data.username} role={data.data.role} />
        </>
      ):
<form onSubmit={handleSubmit}>
  <div className="input-box">
        <p className="">Email:</p>
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
        <p className="warning">Please enter a valid email address.</p>
        )}
        </div>
        <div className="input-box">
        <p className="">Password:</p>
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
        <p className="warning">
          Password must be at least 8 characters long.
        </p>
      )}
      </div>
      <button className="login" type="submit">Submit</button>
    </form>
    }
    </main>
  );
}