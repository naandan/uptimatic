import { NextRequest, NextResponse } from "next/server";

export function middleware(req: NextRequest) {
  const { cookies, nextUrl } = req;
  const pathname = nextUrl.pathname;

  const token = cookies.get("refresh_token")?.value;

  // --- Batas route ---
  const authRoutes = ["/auth/login", "/auth/register"];
  const protectedRoutes = ["/uptime"];
  const verifyRoutes = ["/auth/verify", "/auth/verify-success"];
  const resendRoute = "/auth/resend-verification";

  // --- Kalau belum login ---
  if (!token) {
    // Blokir protected route
    if (protectedRoutes.some((p) => pathname.startsWith(p))) {
      return NextResponse.redirect(new URL("/auth/login", req.url));
    }

    // Blokir resend-verification (harus login)
    if (pathname.startsWith(resendRoute)) {
      return NextResponse.redirect(new URL("/auth/login", req.url));
    }

    // Tapi biarkan akses ke /auth/verify & /auth/verify-success
    return NextResponse.next();
  }

  // --- Decode JWT ---
  let verified = false;
  try {
    const payloadBase64 = token.split(".")[1];
    const decodedPayload = JSON.parse(Buffer.from(payloadBase64, "base64").toString());
    verified = decodedPayload.verified === true;
  } catch {
    verified = false; // token rusak
  }

  // --- Kalau sudah login tapi belum verify ---
  if (!verified && protectedRoutes.some((p) => pathname.startsWith(p))) {
    return NextResponse.redirect(new URL("/auth/resend-verification", req.url));
  }

  // --- Kalau sudah verified, jangan bisa buka login/register lagi ---
  if (verified && authRoutes.some((p) => pathname.startsWith(p))) {
    return NextResponse.redirect(new URL("/uptime", req.url));
  }

  // --- Kalau sudah verified, tapi buka /auth/verify atau /auth/resend-verification, arahkan ke /verify-success ---
  if (verified && [verifyRoutes[0], resendRoute].some((p) => pathname === p)) {
    return NextResponse.redirect(new URL("/auth/verify-success", req.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/uptime/:path*",
    "/auth/login",
    "/auth/register",
    "/auth/resend-verification",
    "/auth/verify",
    "/auth/verify-success",
  ],
};
