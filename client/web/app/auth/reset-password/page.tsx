'use client';

import { useState, useEffect, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Lock } from "lucide-react";
import { authService } from "@/lib/services/auth";
import { toast } from "sonner";
import { getErrorMessage } from "@/utils/helper";

function ResetPasswordForm(){
  const searchParams = useSearchParams();
  const token = searchParams.get("token");
  const router = useRouter();
  
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  
  useEffect(() => {
    if (!token) {
      toast.error("Token tidak valid atau tidak tersedia.");
    }
  }, [token]);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
  
    if (!token) return;
    if (password !== confirmPassword) {
      toast.error("Password dan konfirmasi password tidak cocok.");
      return;
    }
  
    setLoading(true);
    const res = await authService.resetPassword(token, password);
    if (!res.success) {
      toast.error(getErrorMessage(res.error?.code || ""));
      setLoading(false);
      return;
    } else {
      toast.success("Password berhasil diubah. Silakan login dengan password baru.");
      router.replace("/auth/reset-success");
      setLoading(false);
    }
  };
  
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <Lock className="w-12 h-12 text-primary mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          Reset Password
        </h1>
        <p className="text-slate-600 mb-4">
          Masukkan password baru kamu di bawah ini.
        </p>
  
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            type="password"
            placeholder="Password baru"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <Input
            type="password"
            placeholder="Konfirmasi password baru"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
          />
  
          <Button type="submit" disabled={loading || !token} className="w-full">
            {loading ? "Memproses..." : "Reset Password"}
          </Button>
        </form>
      </div>
    </div>
  );
}

export default function ResetPasswordPage() {
  return (
     <Suspense fallback={<div>Loading...</div>}>
       <ResetPasswordForm />
     </Suspense>
  )
}
