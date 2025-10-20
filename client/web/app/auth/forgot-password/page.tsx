"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Mail } from "lucide-react";
import { authService } from "@/lib/services/auth";
import { toast } from "sonner";
import { getErrorMessage } from "@/utils/helper";

export default function ForgotPasswordPage() {
  const [email, setEmail] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    const res = await authService.forgotPassword(email);
    if (!res.success) {
      toast.error(getErrorMessage(res.error?.code || ""));
    } else {
      toast.success("Link reset password telah dikirim ke email kamu.");
    }
    setLoading(false);
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50 px-4">
      <div className="bg-white p-8 rounded-2xl shadow-sm max-w-md w-full text-center border border-slate-200">
        <Mail className="w-12 h-12 text-primary mx-auto mb-4" />
        <h1 className="text-xl font-semibold text-slate-800 mb-2">
          Lupa Password
        </h1>
        <p className="text-slate-600 mb-4">
          Masukkan email kamu di bawah ini, kami akan mengirim link untuk
          mereset password.
        </p>

        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <Input
            type="email"
            placeholder="Email kamu"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />

          <Button type="submit" disabled={loading} className="w-full">
            {loading ? "Mengirim..." : "Kirim Link Reset"}
          </Button>
        </form>
      </div>
    </div>
  );
}
