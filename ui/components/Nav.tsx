import Link from "@/node_modules/next/link";
import React from "react";

const Nav = () => {
  return (
    <nav className="p-4 bg-blue-400 text-white text-xl flex items-center gap-8 justify-center">
      <Link href="/" className="">
        Home
      </Link>
      <Link href="/user" className="">
        Profile
      </Link>
      <Link href="/login" className="">
        Login
      </Link>
    </nav>
  );
};

export default Nav;
