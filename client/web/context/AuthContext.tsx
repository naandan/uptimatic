'use client';
import { authService } from "@/lib/services/auth";
import { createContext, useContext, useState, useEffect } from "react";

interface AuthContextType {
  isLoggedIn: boolean;
  setLoggedIn: (val: boolean) => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const getProfile = async () => {
      const res = await authService.profile();
      if (!res.success) {
        setIsLoading(false);
        return;
      } else {
        setIsLoggedIn(true);
        setIsLoading(false);
      }
    };
    getProfile();
  }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn, setLoggedIn: setIsLoggedIn, isLoading }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
