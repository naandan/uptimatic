import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios";
import { ApiResponse } from "@/types/response";

let isRefreshing = false;
let failedQueue: { resolve: (value?: unknown) => void; reject: (reason?: unknown) => void }[] = [];

const processQueue = (error: AxiosError | null) => {
  failedQueue.forEach((prom) => {
    if (error) prom.reject(error);
    else prom.resolve();
  });
  failedQueue = [];
};

const api: AxiosInstance = axios.create({
  baseURL: "/api/v1",
  withCredentials: true,
});

api.interceptors.response.use(
  (response) => {
    const data = response.data || {};
    const normalized: ApiResponse =  {
      success: true,
      requestId: data.request_id || null,
      data: data.data || null,
      meta: data.meta || null,
      error: null,
    }
    return normalized as unknown as AxiosResponse<ApiResponse>;
  },

  async (error: AxiosError|any) => {
    const originalRequest = error.config as any;
    const res = error.response;
    const status = res?.status;

    if (status === 401 && !originalRequest?._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => api(originalRequest))
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        await api.post("/auth/refresh");
        processQueue(null);
        return api(originalRequest);
      } catch (refreshErr) {
        processQueue(refreshErr as AxiosError);
        failedQueue = [];

        if ((refreshErr as AxiosError).response?.status === 401) {
          console.warn("Refresh token invalid — forcing logout");
          // window.location.href = "/auth/login";
        }

        return Promise.reject(refreshErr);
      } finally {
        isRefreshing = false;
      }
    }

    const payload = res?.data || {};
    const errObj = payload.error || {};

    if (status && status >= 500) {
      console.error("Server Error:", error.message);
    }

    const normalizedError: ApiResponse = {
      success: false,
      requestId: payload.request_id || null,
      data: null,
      error: {
        message: errObj.message || error.message || "Unknown error",
        code: errObj.code || "UNKNOWN_ERROR",
        fields: errObj.fields || null,
        status: status || 0,
      },
    };
    
    // return Promise.reject(error);
    return Promise.resolve(normalizedError as unknown as AxiosResponse<ApiResponse>);
  }
);

export default api;
