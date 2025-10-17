"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { authService } from "@/lib/services/auth";
import { useAuth } from "@/context/AuthContext";
import { toast } from "sonner";
import { getErrorMessage, getValidationErrors } from "@/utils/helper";
import { ErrorInput } from "@/types/response";
import ErrorInputMessage from "./ErrorInputMessage";
import { FcGoogle } from "react-icons/fc";

interface AuthFormProps {
  type: "login" | "register";
}

export const AuthForm = ({ type }: AuthFormProps) => {
  const { setLoggedIn } = useAuth();
  const router = useRouter();
  const [payload, setPayload] = useState({
    email: "",
    password: "",
  });
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<ErrorInput[]>();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    if (type === "register") {
      const res = await authService.register(payload);
      if (!res.success) {
        if (res.error?.code === "VALIDATION_ERROR") {
          const details = getValidationErrors(res.error.fields);
          setErrors(details);
          toast.error(getErrorMessage(res.error?.code || ""));
        } else {
          toast.error(getErrorMessage(res.error?.code || ""));
        }
      }
      setLoading(false);
    }else{
      const res = await authService.login(payload);
      if (!res.success) {
        toast.error(getErrorMessage(res.error?.code || ""));
      }else{
        toast.success("Login berhasil");
        setLoggedIn(true);
        router.push("/uptime");
      }
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-b from-slate-50 to-slate-100 p-6">
      <motion.div
        initial={{ opacity: 0, y: 16 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.4 }}
        className="w-full max-w-sm"
      >
        <Card className="shadow-lg rounded-2xl border border-slate-200">
          <CardHeader className="text-center space-y-2">
            <CardTitle className="text-3xl font-semibold tracking-tight">
              {type === "login" ? "Masuk" : "Daftar"}
            </CardTitle>
            <CardDescription>
              {type === "login"
                ? "Masuk ke akun Uptimatic untuk memantau uptime dan performa layanan Anda."
                : "Daftar sekarang dan mulai pantau uptime website Anda secara real-time."}
            </CardDescription>
          </CardHeader>

          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-5">
              <div className="space-y-2">
                <Label htmlFor="email">Alamat Email</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="contoh@domain.com"
                  autoComplete="email"
                  value={payload.email}
                  onChange={(e) =>
                    setPayload({ ...payload, email: e.target.value })
                  }
                  required
                />
                <ErrorInputMessage errors={errors} field="email" />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">Kata Sandi</Label>
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  autoComplete="current-password"
                  value={payload.password}
                  onChange={(e) =>
                    setPayload({ ...payload, password: e.target.value })
                  }
                  required
                />
                <ErrorInputMessage errors={errors} field="password" />
                <div className="flex items-center justify-end">
                {type === "login" && (
                        <Link
                        href="/auth/forgot-password"
                        className="text-xs text-primary hover:underline"
                        >
                        Lupa kata sandi?
                        </Link>
                    )}
                </div>
              </div>


              <Button
                type="submit"
                className="w-full"
                disabled={loading}
              >
                {loading
                  ? "Memproses..."
                  : type === "login"
                  ? "Masuk"
                  : "Daftar"}
              </Button>

              <div className="relative">
                <div className="absolute inset-0 flex items-center">
                  <span className="w-full border-t" />
                </div>
                <div className="relative flex justify-center text-xs uppercase">
                  <span className="bg-white px-2 text-slate-500">
                    Atau masuk dengan
                  </span>
                </div>
              </div>

              <div className="flex gap-4">
                <Button variant="outline" onClick={() => window.location.replace("/api/v1/auth/google/login")} type="button" className="w-full">
                  <FcGoogle size={24} />
                  Google
                </Button>
              </div>
            </form>
          </CardContent>

          <CardFooter className="text-sm text-muted-foreground flex flex-col">
            {type === "login" ? (
              <p>
                Belum punya akun?{" "}
                <Link
                  href="/auth/register"
                  className="text-primary hover:underline"
                >
                  Daftar
                </Link>
              </p>
            ) : (
              <p>
                Sudah punya akun?{" "}
                <Link
                  href="/auth/login"
                  className="text-primary hover:underline"
                >
                  Masuk
                </Link>
              </p>
            )}
          </CardFooter>
        </Card>

        <div className="mt-4 text-center text-xs text-slate-400">
          © {new Date().getFullYear()} Uptimatic. Semua hak dilindungi.
        </div>
      </motion.div>
    </div>
  );
}
