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
    }
};
