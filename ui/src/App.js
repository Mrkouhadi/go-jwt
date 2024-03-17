import { useEffect, useState } from "react";

export default function Home() {
  const [data, setData] = useState();

  async function RefreshAT() {
      try {
        const res = await fetch("http://localhost:8080/refresh-token", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
        });
        if (!res.ok) {
          throw new Error("Failed to fetch data");
        }
        const jsonData = await res.json();
        setData(jsonData);
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
            username: "bryan",
            email: "bryan@hotmail.com",
          }),
          headers: {
            "Content-Type": "application/json",
          },
        });
        if (!res.ok) {
          throw new Error("Failed to fetch data");
        }
        const jsonData = await res.json();
        setData(jsonData);
      } catch (error) {
        console.error("Error fetching data:", error);
      }
    }
    getData();
  },[]);

  return (
    <main className="flex h-screen flex-col items-center gap-12 p-12 bg-pink-400 text-black">
      {data && (
        <>
          <p>Access Token: {data.access_token}</p>
          <p>Refresh Token: {data.refresh_token}</p>
        </>
      )}

      <button className="px-4 py-2" onClick={RefreshAT}>Rfersh the token</button>
    </main>
  );
}
