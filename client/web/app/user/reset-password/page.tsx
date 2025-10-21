"use client";

import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { useState } from "react";
import Link from "next/link";
import { userService } from "@/lib/services/user";
import { toast } from "sonner";

export default function ProfilePage() {
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      toast.error("Password dan konfirmasi password tidak cocok.");
      return;
    }

    const res = await userService.changePassword(currentPassword, newPassword);
    if (!res.success) {
      toast.error("Terjadi kesalahan saat mengubah password.");
    } else {
      toast.success("Password berhasil diubah.");
      setCurrentPassword("");
      setNewPassword("");
      setConfirmPassword("");
    }
  };

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Reset Kata Sandi</h1>

      {/* Form Reset Password */}
      <form className="space-y-5 max-w-lg" onSubmit={handleSubmit}>
        <div className="space-y-2">
          <Label htmlFor="currentPassword">Password Lama</Label>
          <Input
            id="currentPassword"
            type="password"
            placeholder="Masukkan password lama"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
          />
          <p className="text-sm text-primary text-right">
            <Link href="/auth/forgot-password" className="hover:underline">
              Lupa Password?
            </Link>
          </p>
        </div>

        <div className="space-y-2">
          <Label htmlFor="newPassword">Password Baru</Label>
          <Input
            id="newPassword"
            type="password"
            placeholder="Masukkan password baru"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="confirmPassword">Konfirmasi Password Baru</Label>
          <Input
            id="confirmPassword"
            type="password"
            placeholder="Masukkan konfirmasi password baru"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
          />
        </div>

        <Button type="submit" className="w-full md:w-auto mt-1">
          Simpan
        </Button>
      </form>
    </div>
  );
}
