import api from "../api";

export type URLRequest = {
    label: string,
    url: string,
    active: boolean
}

export const urlService = {
    create: async (data: URLRequest) => {
        const res = await api.post("/urls", data);
        return res.data;
    },

    update: async (id: number, data: URLRequest) => {
        const res = await api.put(`/urls/${id}`, data);
        return res.data;
    },

    delete: async (id: number) => {
        const res = await api.delete(`/urls/${id}`);
        return res.data;
    },

    get: async (id: number) => {
        const res = await api.get(`/urls/${id}`);
        return res.data;
    },

    list: async (params?: string) => {
        const res = await api.get(`/urls?${params}`);
        return res.data;
    },

    stats: async (id: number, mode: "day" | "month", offset: number) => {
        const res = await api.get(`/urls/${id}/stats?mode=${mode}&offset=${offset}`);
        return res.data;
    }
}