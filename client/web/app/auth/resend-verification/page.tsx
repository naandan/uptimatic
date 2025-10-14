"use client";

import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Mail } from "lucide-react";
import { useRouter } from "next/navigation";
import { authService } from "@/lib/services/auth";
import { toast } from "sonner";

export default function VerifyEmailPage() {
  const router = useRouter();

  const [countdown, setCountdown] = useState(0);
  const [loading, setLoading] = useState(true);
  const [hasStarted, setHasStarted] = useState(false);

  useEffect(() => {
    if (countdown <= 0) return;
    const timer = setInterval(() => setCountdown((prev) => prev - 1), 1000);
    return () => clearInterval(timer);
  }, [countdown]);

  useEffect(() => {
    const checkLogin = async () => {
      try {
        setLoading(true);
        await authService.profile();

        const res = await authService.resendVerificationEmailTTL();
        if (res.data.ttl > 0) {
          setCountdown(res.data.ttl);
          setHasStarted(true);
        }
      } catch {
        router.replace("/auth/login");
      }finally {
        setLoading(false);
      } 
    };
    checkLogin();
  }, [router]);

  const handleSend = async () => {
    try {
      setLoading(true);
      const res = await authService.resendVerificationEmail();
      toast.success("Email verifikasi telah dikirim.");
      setCountdown(res.data.ttl);
      setHasStarted(true);
    } catch (err: any) {
      toast.error(err.message || "Gagal mengirim email verifikasi.");
    } finally {
      setLoading(false);
    }
  };

  const handleRefresh = async () => {
    try {
      await authService.refresh();
      window.location.reload();
    } catch (err) {
      console.error(err);
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

        <Button onClick={handleRefresh} className="w-full">
          Refresh
        </Button>
      </div>
    </div>
  );
}
