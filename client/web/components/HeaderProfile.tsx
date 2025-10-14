"use client";

import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { useAuth } from "@/context/AuthContext";
import { authService } from "@/lib/services/auth";
import { User } from "lucide-react";
import { useRouter } from "next/navigation";

export default function HeaderProfile() {
  const { isLoggedIn, setLoggedIn } = useAuth();
  const router = useRouter();
  

  const handleLogout = async () => {
    try{
      await authService.logout();
      setLoggedIn(false);
      router.push("/auth/login");
    }catch(err){
      console.error(err);
    }
  };

  if (!isLoggedIn) return null;

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="w-10 h-10 p-0 rounded-full">
          {/* Bisa diganti dengan foto user */}
          <User className="w-6 h-6" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-32">
        <DropdownMenuItem onClick={handleLogout}>Logout</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
