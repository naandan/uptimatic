import { ApiError, ValidationFieldError } from "@/types/response";

const VALIDATION: Record<string, string | ((param?: string) => string)> = {
  REQUIRED: "Tidak boleh kosong.",
  INVALID_TYPE: "Data tidak valid.",
  INVALID_FORMAT: "Format tidak valid.",
  MIN_LENGTH: (min?: string) => `Minimal ${min} karakter.`,
  MAX_LENGTH: (max?: string) => `Maksimal ${max} karakter.`,
  ENUM_VALUE: "Nilai tidak valid.",
  MISMATCH: "Tidak cocok.",
  UNIQUE: "Sudah digunakan.",
  MIN_VALUE: (min?: string) => `Minimal ${min}.`,
  MAX_VALUE: (max?: string) => `Maksimal ${max}.`,
};


const MESSAGES: Record<string, string> = {
  VALIDATION_ERROR: "Input tidak valid.",

  UNAUTHORIZED: "Anda belum login.",
  INVALID_CREDENTIALS: "Email atau password salah.",
  INVALID_TOKEN: "Token tidak valid.",
  FORBIDDEN_ACTION: "Aksi tidak diizinkan.",
  ACCOUNT_LOCKED: "Akun kamu terkunci.",
  TOO_MANY_REQUESTS: "Terlalu banyak permintaan.",

  NOT_FOUND: "Data tidak ditemukan.",
  CONFLICT: `Data sudah digunakan.`,

  INTERNAL_ERROR: "Terjadi kesalahan yang tidak diketahui.",
  SERVICE_UNAVAILABLE: "Layanan tidak tersedia.",
  TIMEOUT: "Koneksi terputus.",
};


export function getValidationErrors(fields?: ApiError["fields"]) {
  if (!fields) return [];

  return Object.entries(fields).map(([field, errors]) => {
    const reasons = errors.map((err: ValidationFieldError) => {
      const msgTemplate = VALIDATION[err.code];
      if (typeof msgTemplate === "function") {
        return msgTemplate(err.param);
      }
      return msgTemplate || err.code;
    });

    return { field, reasons };
  });
}

export function getErrorMessage(code: string): string {
  return MESSAGES[code] || "Terjadi kesalahan yang tidak diketahui.";
}
