"use client";

import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { ServerCrash } from "lucide-react";

export default function InternalErrorPage() {
  const router = useRouter();

  return (
    <div className="min-h-[77vh] flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <ServerCrash className="w-12 h-12 text-red-600 mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          Terjadi Kesalahan
        </h1>
        <p className="text-slate-600 mb-4">
          Maaf, terjadi kesalahan pada server. Silakan coba lagi nanti.
        </p>
        <Button onClick={() => router.refresh()} className="w-full">
          Muat Ulang Halaman
        </Button>
      </div>
    </div>
  );
}
