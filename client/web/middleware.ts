import { NextRequest, NextResponse } from "next/server";

export function middleware(req: NextRequest) {
  const { cookies, nextUrl } = req;
  const pathname = nextUrl.pathname;

  const token = cookies.get("refresh_token")?.value;

  // --- Batas route ---
  const authRoutes = ["/auth/login", "/auth/register"];
  const verifyRoutes = ["/auth/verify", "/auth/verify-success"];
  const protectedRoutes = ["/uptime"];

  // --- Kalau belum login ---
  if (!token) {
    if (protectedRoutes.some((p) => pathname.startsWith(p))) {
      return NextResponse.redirect(new URL("/auth/login", req.url));
    }
    return NextResponse.next();
  }

  // --- Decode JWT ---
  let verified = false;
  try {
    const payloadBase64 = token.split(".")[1];
    const decodedPayload = JSON.parse(Buffer.from(payloadBase64, "base64").toString());
    verified = decodedPayload.verified === true;
  } catch {
    // Jika token rusak, anggap belum verified
    verified = false;
  }

  // --- Kalau sudah login tapi belum verify ---
  if (!verified && protectedRoutes.some((p) => pathname.startsWith(p))) {
    return NextResponse.redirect(new URL("/auth/verify", req.url));
  }

  // --- Kalau sudah verified, jangan bisa buka /auth/verify & /auth/verify-success ---
  if (verified && verifyRoutes.some((p) => pathname.startsWith(p))) {
    return NextResponse.redirect(new URL("/uptime", req.url));
  }

  // --- Kalau sudah login, jangan bisa buka login/register lagi ---
  if (token && authRoutes.some((p) => pathname.startsWith(p))) {
    return NextResponse.redirect(new URL("/uptime", req.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/uptime/:path*",
    "/auth/login",
    "/auth/register",
    "/auth/verify",
    "/auth/verify-success",
  ],
};
