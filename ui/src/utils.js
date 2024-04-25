
// logging in the user
export  async function Login(email,password,setData) {
 try {
   const res = await fetch("/login", {
     method: "POST",
     body: JSON.stringify({
       email: email,
       password: password,
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

// logging out the user
export async function Logout(setData){
    try {
      const res = await fetch("/logout",{method:"GET"});
      if (!res.ok) {
        throw new Error("Failed to fetch data");
      }
      setData(null);
    } catch (error) {
      console.error("Error refreshing token:", error);
    }
  }