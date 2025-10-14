'use client';

import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { CheckCircle } from "lucide-react";

export default function ResetSuccessPage() {
  const router = useRouter();

  const handleLogin = () => {
    router.replace("/auth/login");
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <CheckCircle className="w-12 h-12 text-green-600 mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          Reset Password Berhasil
        </h1>
        <p className="text-slate-600 mb-4">
          Password kamu sudah berhasil diubah. Silakan login menggunakan password baru.
        </p>

        <Button onClick={handleLogin} className="w-full">
          Login Sekarang
        </Button>
      </div>
    </div>
  );
}
