"use client";

import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { FileWarning } from "lucide-react";

export default function NotFoundPage() {
  const router = useRouter();

  return (
    <div className="min-h-[77vh] flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <FileWarning className="w-12 h-12 text-yellow-500 mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          Halaman Tidak Ditemukan
        </h1>
        <p className="text-slate-600 mb-4">
          Maaf, halaman yang kamu cari tidak ada.
        </p>
        <Button onClick={() => router.replace("/")} className="w-full">
          Kembali ke Beranda
        </Button>
      </div>
    </div>
  );
}
