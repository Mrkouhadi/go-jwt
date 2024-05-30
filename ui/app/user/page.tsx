"use client";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function User() {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);

  const router = useRouter();

  // fetch user data
  useEffect(() => {
    setLoading(true);
    fetch("http://localhost:8080/protected/profile", {
      method: "GET",
      credentials: "include",
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(
            "Failed to fetch profile data. Status: " + response.status
          );
        }
        return response.json();
      })
      .then((data) => {
        setData(data);
      })
      .catch((error) => console.error("Error fetching profile data:", error))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  // logout handler
  const logout = async () => {
    try {
      const response = await fetch("http://localhost:8080/logout", {
        method: "GET",
      });
      if (!response.ok) {
        throw new Error(
          "Logout request failed with status: " + response.status
        );
      }
      const contentType = response.headers.get("content-type");
      if (contentType && contentType.includes("application/json")) {
        const data = await response.json();
        setData(null);
        router.push("/login");
      } else {
        setData(null);
        router.push("/login");
      }
    } catch (error) {
      console.error("Error logging out: ", error);
    }
  };

  return (
    <div>
      <h1>Fetched Data</h1>
      <pre>{data && JSON.stringify(data, null, 2)}</pre>
      <button onClick={logout} className="bg-blue-400 px-8 py-4 text-xl">
        Log out
      </button>
    </div>
  );
}
