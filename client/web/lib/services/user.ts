import { ApiResponse } from "@/types/response";
import api from "../api";
import {
  UpdateFotoResponse,
  UploadURLResponse,
  UserResponse,
} from "@/types/user";

export const userService = {
  me: async (): Promise<ApiResponse<UserResponse>> => {
    return await api.get("/users/me");
  },

  update: async (
    name: string,
    email: string,
  ): Promise<ApiResponse<UserResponse>> => {
    return await api.put("/users", { name, email });
  },

  changePassword: async (
    old_password: string,
    new_password: string,
  ): Promise<ApiResponse<null>> => {
    return await api.put("/users/change-password", {
      old_password,
      new_password,
    });
  },

  uploadURL: async (
    file_name: string,
    content_type: string,
  ): Promise<ApiResponse<UploadURLResponse>> => {
    return await api.post("/users/upload-url", { file_name, content_type });
  },

  updateFoto: async (
    file_name: string,
  ): Promise<ApiResponse<UpdateFotoResponse>> => {
    return await api.put("/users/update-foto", { file_name });
  },
};
