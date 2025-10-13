import api from "../api";


type AuthRequest = {
    email: string;
    password: string;
}

export const authService = {
    register: async (data: AuthRequest) => {
        const res = await api.post("/auth/register", data);
        return res.data;
    },

    login: async (data: AuthRequest) => {
        const res = await api.post("/auth/login", data);
        return res.data;
    },

    logout: async () => {
        const res = await api.post("/auth/logout");
        return res.data;
    },

    refresh: async () => {
        const res = await api.post("/auth/refresh");
        return res.data;
    },

    profile: async () => {
        const res = await api.get("/profile");
        return res.data;
    },

    verify: async (token: string) => {
        const res = await api.get(`/auth/verify?token=${token}`);
        return res.data;
    },

    resendVerificationEmail: async () => {
        const res = await api.post(`/auth/resend-verification`);
        return res.data;
    }
};
