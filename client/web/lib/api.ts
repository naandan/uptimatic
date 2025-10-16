import axios, { AxiosError, AxiosInstance, AxiosResponse } from "axios";
import { ApiResponse } from "@/types/response";

let isRefreshing = false;
let failedQueue: { resolve: (value?: any) => void; reject: (err: any) => void }[] = [];

const processQueue = (error: any | null) => {
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
    const normalized: ApiResponse = {
      success: true,
      requestId: data.request_id || null,
      data: data.data || null,
      meta: data.meta || null,
      error: null,
    };
    return normalized as unknown as AxiosResponse<ApiResponse>;
  },

  async (error: AxiosError | any) => {
    const originalRequest = error.config;
    const res = error.response;
    const status = res?.status;

    if (status === 401 && !originalRequest?._retry) {
      originalRequest._retry = true;

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => api(originalRequest))
          .catch(() =>
            Promise.resolve({ success: false, error: { message: "Unauthorized", code: "UNAUTHORIZED" } })
          );
      }

      isRefreshing = true;

      try {
        await axios.post("/api/v1/auth/refresh", {}, { withCredentials: true });
        processQueue(null);
        return api(originalRequest);
      } catch (refreshErr) {
        processQueue(refreshErr);
        failedQueue = [];
        return Promise.resolve({
          success: false,
          error: { message: "Unauthorized", code: "UNAUTHORIZED" },
        });
      } finally {
        isRefreshing = false;
      }
    }

    const payload = res?.data || {};
    const errObj = payload.error || {};

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

    return Promise.resolve(normalizedError);
  }
);

export default api;
