import api from "../api";

type URLRequest = {
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

    list: async () => {
        const res = await api.get("/urls");
        return res.data;
    },

    stats: async (id: number) => {
        const res = await api.get(`/urls/${id}/stats`);
        return res.data;
    }
}