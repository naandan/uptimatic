export type AuthRequest = {
  email: string;
  password: string;
};

export type TokenResponse = {
  access_token: string;
  refresh_token: string;
};

export type TTLResponse = {
  ttl: number;
};
