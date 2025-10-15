"use client";

import { Suspense, useEffect, useState } from "react";
import { Mail } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { authService } from "@/lib/services/auth";

function VerifyEmailForm(){

  const searchParams = useSearchParams();
  const token = searchParams.get("token");
  const router = useRouter();
  
  const [message, setMessage] = useState("Memverifikasi email kamu...");
  const [isError, setIsError] = useState(false);
  
  useEffect(() => {
    if (!token) {
      router.replace("/auth/login");
      return;
    }
  
    const verify = async () => {
      const res = await authService.verify(token);
      if (!res.success) {
        console.error("Verifikasi gagal:", res.error);
        setIsError(true);
        setMessage("Token tidak valid atau sudah kedaluwarsa.");
        return;
      } else {
        const res = await authService.profile();
        if (res.success){
          await authService.refresh();
        }
        router.replace("/auth/verify-success");
      }
    };
  
    verify();
  }, [token, router]);
  
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <Mail className="w-12 h-12 text-primary mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          {isError ? "Verifikasi Gagal" : "Memverifikasi Email..."}
        </h1>
        <p
          className={`text-slate-600 ${
            isError ? "text-red-500" : "text-slate-600"
          }`}
        >
          {message}
        </p>
      </div>
    </div>
  );
}

export default function VerifyEmailPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <VerifyEmailForm />
    </Suspense>
  );
}
