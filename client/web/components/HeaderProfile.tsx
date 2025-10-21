"use client";

import { useAuth } from "@/context/AuthContext";
import { authService } from "@/lib/services/auth";
import { useRouter } from "next/navigation";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
  DropdownMenuLabel,
} from "@/components/ui/dropdown-menu";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { LogOut, Settings, User } from "lucide-react";

export default function HeaderProfile() {
  const { isLoggedIn, isLoading, user, setIsLoggedIn } = useAuth();
  const router = useRouter();

  const handleLogout = async () => {
    const res = await authService.logout();
    if (!res.success) return;
    setIsLoggedIn(false);
    router.push("/auth/login");
  };

  if (isLoading || !isLoggedIn || !user) return null;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          variant="ghost"
          className="w-10 h-10 p-0 rounded-full"
          aria-label="User menu"
        >
          <Avatar className="w-9 h-9">
            <AvatarImage
              src={
                user.profile ||
                `https://placehold.co/100x100?text=${user?.name?.slice(0, 1)}`
              }
            />
            <AvatarFallback>
              {user.name
                ? user.name
                    .split(" ")
                    .map((n) => n[0])
                    .join("")
                    .slice(0, 2)
                    .toUpperCase()
                : "U"}
            </AvatarFallback>
          </Avatar>
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="end" className="w-56 p-2">
        {/* Header user info */}
        <DropdownMenuLabel className="flex flex-col items-start">
          <p className="font-semibold text-sm">{user.name}</p>
          <p className="text-xs text-gray-500 truncate">{user.email}</p>
        </DropdownMenuLabel>

        <DropdownMenuSeparator />

        {/* Menu items */}
        <DropdownMenuItem onClick={() => router.push("/user/profile")}>
          <User className="w-4 h-4 mr-2" /> Profil Saya
        </DropdownMenuItem>

        <DropdownMenuItem onClick={() => router.push("/user/settings")}>
          <Settings className="w-4 h-4 mr-2" /> Pengaturan
        </DropdownMenuItem>

        <DropdownMenuSeparator />

        <DropdownMenuItem
          onClick={handleLogout}
          className="text-red-600 focus:text-red-700"
        >
          <LogOut className="w-4 h-4 mr-2" /> Logout
        </DropdownMenuItem>

        {/* Optional: footer */}
        <div className="mt-2 text-[10px] text-gray-400 text-center select-none">
          Versi 1.0.0
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
