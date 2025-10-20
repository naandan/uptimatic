"use client";

import { usePathname } from "next/navigation";
import { ReactNode } from "react";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { Toaster } from "@/components/ui/sonner";
import { useAuth } from "@/context/AuthContext";
import Image from "next/image";

export const ClientLayout = ({ children }: { children: ReactNode }) => {
  const { isLoading } = useAuth();
  const pathname = usePathname();
  const hideHeaderFooter =
    pathname?.startsWith("/auth") || pathname?.startsWith("/admin");

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Image
          src="/icon.png"
          alt="Loading"
          width={100}
          height={100}
          className="animate-ping"
        ></Image>
      </div>
    );
  }
  return (
    <>
      <Toaster />
      {!hideHeaderFooter && <Header />}
      <main>{children}</main>
      {!hideHeaderFooter && <Footer />}
    </>
  );
};
