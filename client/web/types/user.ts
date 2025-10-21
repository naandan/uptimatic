export type UserResponse = {
  id: number;
  name: string;
  email: string;
  verified: boolean;
  profile: string;
  created_at: string;
};

export type UploadURLResponse = {
  file_name: string;
  presigned_url: string;
};

export type UpdateFotoResponse = {
  url: string;
};
