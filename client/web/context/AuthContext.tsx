'use client';
import { authService } from "@/lib/services/auth";
import { createContext, useContext, useState, useEffect } from "react";

interface AuthContextType {
  isLoggedIn: boolean;
  setLoggedIn: (val: boolean) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const getProfile = async () => {
      const res = await authService.profile();
      if (!res.success) {
        setIsLoggedIn(false);
        return;
      } else {
        setIsLoggedIn(true);
      }
    }
    getProfile();
  }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn, setLoggedIn: setIsLoggedIn }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
