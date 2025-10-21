"use client";

import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useState } from "react";

export default function ProfilePage() {
  const [securityEnabled, setSecurityEnabled] = useState(false);

  const toggleSecurityEnabled = () => setSecurityEnabled((prev) => !prev);

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Keamanan</h1>

      <div className="flex items-center space-x-2">
        <Switch
          id="airplane-mode"
          onChange={toggleSecurityEnabled}
          checked={securityEnabled}
        />
        <Label htmlFor="airplane-mode">2 Factor Authentication</Label>
      </div>
      <p className="text-sm text-gray-500">
        2FA adalah teknologi untuk memastikan bahwa pengguna memiliki akses ke
        akun mereka dengan menggunakan dua faktor seperti kode OTP (One-Time
        Password) atau token yang diberikan oleh penyedia layanan.
      </p>
    </div>
  );
}
