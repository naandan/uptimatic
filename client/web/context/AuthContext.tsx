"use client";
import { userService } from "@/lib/services/user";
import { UserResponse } from "@/types/user";
import { createContext, useContext, useState, useEffect } from "react";

interface AuthContextType {
  isLoggedIn: boolean;
  setIsLoggedIn: (val: boolean) => void;
  isLoading: boolean;
  user: UserResponse | null;
  setUser: (val: UserResponse | null) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [user, setUser] = useState<UserResponse | null>(null);

  useEffect(() => {
    const getProfile = async () => {
      const res = await userService.me();
      if (res.success) {
        setUser(res.data);
        setIsLoggedIn(true);
      } else {
        setIsLoggedIn(false);
      }

      setIsLoading(false);
    };

    getProfile();
  }, []);

  return (
    <AuthContext.Provider
      value={{ isLoggedIn, setIsLoggedIn, isLoading, user, setUser }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
