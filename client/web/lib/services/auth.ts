import { AuthRequest, TokenResponse, TTLResponse, UserResponse } from "@/types/auth";
import api from "../api";
import { ApiResponse } from "@/types/response";

export const authService = {
    register: async (data: AuthRequest): Promise<ApiResponse<UserResponse>> => {
        return await api.post("/auth/register", data);
    },

    login: async (data: AuthRequest): Promise<ApiResponse<TokenResponse>> => {
        return await api.post("/auth/login", data);
    },

    logout: async (): Promise<ApiResponse<null>> => {
        return await api.post("/auth/logout");
    },

    refresh: async (): Promise<ApiResponse<TokenResponse>> => {
        return await api.post("/auth/refresh");
    },

    profile: async (): Promise<ApiResponse<UserResponse>>  => {
        return await api.get("/profile");
    },

    verify: async (token: string): Promise<ApiResponse<null>> => {
        return await api.get(`/auth/verify?token=${token}`);
    },

    resendVerificationEmail: async (): Promise<ApiResponse<TTLResponse>> => {
        return await api.post(`/auth/resend-verification`);
    },

    resendVerificationEmailTTL: async (): Promise<ApiResponse<TTLResponse>> => {
        return await api.get(`/auth/resend-verification-ttl`);
    },

    forgotPassword: async (email: string): Promise<ApiResponse<null>> => {
        return await api.post(`/auth/forgot-password`, { email });
    },

    resetPassword: async (token: string, password: string): Promise<ApiResponse<null>> => {
        return await api.post(`/auth/reset-password`, { token, password });
    },
};
