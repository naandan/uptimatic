import { NextRequest, NextResponse } from "next/server";

export function middleware(req: NextRequest) {
  const { cookies, nextUrl } = req;
  const pathname = nextUrl.pathname;

  const authRoutes = ["/auth/login", "/auth/register"];
  const protectedRoutes = ["/uptime"]; // bisa ditambah kalau perlu

  const token = cookies.get("refresh_token")?.value;

  // Jika user sudah login dan mencoba ke /auth/login atau /auth/register → redirect
  if (token && authRoutes.some((path) => pathname.startsWith(path))) {
    return NextResponse.redirect(new URL("/uptime", req.url));
  }

  // Jika user belum login dan mencoba akses halaman protected → redirect ke login
  if (!token && protectedRoutes.some((path) => pathname.startsWith(path))) {
    return NextResponse.redirect(new URL("/auth/login", req.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/uptime/:path*",
    "/auth/login",
    "/auth/register",
  ],
};
