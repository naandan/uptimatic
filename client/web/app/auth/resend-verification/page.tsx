"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Mail } from "lucide-react";
import { useRouter } from "next/navigation";
import { authService } from "@/lib/services/auts";

export default function VerifyEmailPage() {
  const router = useRouter();

  const [countdown, setCountdown] = useState(0);
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [hasStarted, setHasStarted] = useState(false);

  // Jalankan countdown setiap detik
  useEffect(() => {
    if (countdown <= 0) return;
    const timer = setInterval(() => setCountdown((prev) => prev - 1), 1000);
    return () => clearInterval(timer);
  }, [countdown]);

  // Pastikan user sudah login sebelum bisa mengirim email verifikasi
  useEffect(() => {
    const checkLogin = async () => {
      try {
        await authService.profile(); // Jika tidak login, akan error
      } catch {
        router.replace("/auth/login");
      }
    };
    checkLogin();
  }, [router]);

  const handleSend = async () => {
    try {
      setLoading(true);
      setMessage("");
      const res = await authService.resendVerificationEmail();
      console.log(res);
      setMessage("Email verifikasi telah dikirim.");
      setCountdown(60);
      setHasStarted(true);
    } catch (err: any) {
      setMessage(err.message || "Gagal mengirim email verifikasi.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <Mail className="w-12 h-12 text-primary mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          Verifikasi Email Kamu
        </h1>
        <p className="text-slate-600 mb-4">
          Klik tombol di bawah ini untuk mengirim email verifikasi.  
          Setelah terkirim, kamu bisa mengirim ulang setelah 60 detik.
        </p>

        <Button
          onClick={handleSend}
          disabled={loading || countdown > 0}
          className="w-full mb-2"
        >
          {loading
            ? "Mengirim..."
            : countdown > 0
            ? `Kirim ulang dalam ${countdown}s`
            : hasStarted
            ? "Kirim Ulang Email"
            : "Kirim Email Verifikasi"}
        </Button>

        {message && (
          <p
            className={`text-sm mt-2 ${
              message.toLowerCase().includes("gagal")
                ? "text-red-500"
                : "text-green-600"
            }`}
          >
            {message}
          </p>
        )}
      </div>
    </div>
  );
}
