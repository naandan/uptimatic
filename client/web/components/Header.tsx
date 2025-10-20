"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/context/AuthContext";
import HeaderProfile from "./HeaderProfile";

export const Header = () => {
  const { isLoggedIn } = useAuth();
  return (
    <header className="flex justify-between items-center py-4 px-6 md:px-12 bg-white/70 backdrop-blur-md border-b border-slate-200 sticky top-0 z-50">
      <Link href="/" className="text-2xl font-semibold text-primary">
        Uptimatic
      </Link>
      {/* <nav className="hidden md:flex space-x-6 text-sm font-medium text-slate-700">
        <Link href="#features" className="hover:text-primary">Fitur</Link>
        <Link href="#pricing" className="hover:text-primary">Harga</Link>
        <Link href="#about" className="hover:text-primary">Tentang</Link>
      </nav> */}
      <div className="space-x-2">
        {isLoggedIn ? (
          <HeaderProfile />
        ) : (
          <>
            <Button variant="ghost" asChild>
              <Link href="/auth/login">Masuk</Link>
            </Button>
            <Button asChild>
              <Link href="/auth/register">Daftar</Link>
            </Button>
          </>
        )}
      </div>
    </header>
  );
};
