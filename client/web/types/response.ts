export interface ErrorInput { 
  field: string; 
  reasons: string[]
}

export interface ValidationFieldError {
  code: string;
  param?: string;
}

export interface ApiErrorFields {
  [field: string]: ValidationFieldError[];
}

export interface ApiError {
  message: string;
  code: string;
  fields?: ApiErrorFields | null;
  status?: number;
}

export interface ApiMeta {
  total: number;
  limit: number;
  page: number;
  total_page: number;
}

export interface ApiResponse<T = any> {
  success: boolean;
  requestId: string | null;
  data: T | null;
  meta?: ApiMeta | null;
  error: ApiError | null;
}
