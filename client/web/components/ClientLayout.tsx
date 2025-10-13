"use client";

import { usePathname } from "next/navigation";
import { ReactNode } from "react";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { Toaster } from "@/components/ui/sonner";

export const ClientLayout = ({ children }: { children: ReactNode }) => {
  const pathname = usePathname();
  const hideHeaderFooter = pathname?.startsWith("/auth") || pathname?.startsWith("/admin");

  return (
    <>
      <Toaster />
      {!hideHeaderFooter && <Header />}
      <main>{children}</main>
      {!hideHeaderFooter && <Footer />}
    </>
  );
}
