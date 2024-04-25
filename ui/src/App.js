import { useState } from "react";
import { Login, Logout } from "./utils";

export default function Home() {
  const [data, setData] = useState();

  /*
  error : false
  message : "Loggedin user of email email@example.com and Username: bryan kouhadi"
  data:{userid: 'id_ftdrs42671hdcn', username: 'bryan kouhadi', email: 'email@example.com', role: 'user'}
  */
  return (
    <main className="flex h-screen flex-col items-center gap-12 p-12 bg-pink-400 text-black">
      {data ? (
        <>
          <p>{data.message}</p>
          <p>email:{data.data.email}</p>
          <p>username:{data.data.username}</p>
          <p>role:{data.data.role}</p>
        </>
      ):
      <p>No tokens stored yet !</p>
    }
      <button className="px-4 py-2" onClick={()=>Login("email@example.com","123123",setData)}>Log in</button>
      <button className="px-4 py-2 bg-red-400 text-white" onClick={()=>Logout(setData)}>Logout</button>
    </main>
  );
}