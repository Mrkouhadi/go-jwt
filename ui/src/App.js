import { useEffect, useState } from "react";

export default function Home() {
  const [data, setData] = useState();

    const Logout = async()=>{
      try {
        const res = await fetch("http://localhost:8080/logout",{method:"GET"});
        if (!res.ok) {
          throw new Error("Failed to fetch data");
        }
        setData([]);
      } catch (error) {
        console.error("Error refreshing token:", error);
      }
    }
    
  useEffect(() => {
    async function getData() {
      try {
        const res = await fetch("http://localhost:8080/generate-tokens", {
          method: "POST",
          body: JSON.stringify({
            username: "zineb",
            email: "zineb@zineb.com",
          }),
          headers: {
            "Content-Type": "application/json",
          },
        });
        if (!res.ok) {
          throw new Error("Failed to fetch data");
        }
        const jsonData = await res.json();
        setData(jsonData.data);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    }
    getData();
  },[]);


  const refreshAccessToken= async()=>{
    fetch(
      "http://localhost:8080/refresh-token",{
        method:"GET",
        headers:{
          "Content-Type": "application/json",
        }
      }
    )
      .then((response) => response.json())
      .then((data) => {
        setData(data.data)
        console.log(data)
      })
      .catch(function (err) {
        console.log(
          "Unable to fetch -",
          err
        );
      });
  }
  return (
    <main className="flex h-screen flex-col items-center gap-12 p-12 bg-pink-400 text-black">
      {data && (
        <>
          <p>Access Token: {data.access_token}</p>
          <p>Refresh Token: {data.refresh_token}</p>
        </>
      )}
      <button className="px-4 py-2" onClick={refreshAccessToken}>Rfersh the token</button>
      <button className="px-4 py-2 bg-red-400 text-white" onClick={Logout}>Logout</button>
    </main>
  );
}