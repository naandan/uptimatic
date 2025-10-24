import api from "@/lib/api";
import { ApiResponse } from "@/types/response";
import { URLRequest, URLResponse, URLStats } from "@/types/url";

export const urlService = {
  create: async (
    data: Partial<URLRequest>,
  ): Promise<ApiResponse<URLResponse>> => {
    return await api.post("/urls", data);
  },

  update: async (
    id: string | undefined,
    data: Partial<URLRequest>,
  ): Promise<ApiResponse<URLResponse>> => {
    return await api.put(`/urls/${id}`, data);
  },

  delete: async (id: string): Promise<ApiResponse<null>> => {
    return await api.delete(`/urls/${id}`);
  },

  get: async (id: string): Promise<ApiResponse<URLResponse>> => {
    return await api.get(`/urls/${id}`);
  },

  list: async (params?: string): Promise<ApiResponse<URLResponse[]>> => {
    return await api.get(`/urls?${params || ""}`);
  },

  stats: async (
    id: string,
    mode: "day" | "month",
    offset: number,
  ): Promise<ApiResponse<URLStats[]>> => {
    return await api.get(`/urls/${id}/stats?mode=${mode}&offset=${offset}`);
  },
};
