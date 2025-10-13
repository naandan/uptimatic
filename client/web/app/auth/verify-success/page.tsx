"use client";

import { CheckCircle2 } from "lucide-react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useEffect } from "react";
import { authService } from "@/lib/services/auts";

export default function VerifySuccessPage() {
  useEffect(() => {
    const refresh = async () => {
      await authService.refresh();
    };
    refresh();
  })
  
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <CheckCircle2 className="w-12 h-12 text-green-500 mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          Email Berhasil Diverifikasi!
        </h1>
        <p className="text-slate-600 mb-6">
          Akun kamu sudah aktif dan siap digunakan.  
          Silakan lanjut ke dashboard untuk mulai menggunakan layanan.
        </p>
        <Link href="/uptime">
          <Button className="w-full">Ke Dashboard</Button>
        </Link>
      </div>
    </div>
  );
}
