import { NextRequest, NextResponse } from "next/server";

export function middleware(req: NextRequest) {
  const { cookies, nextUrl } = req;
  const pathname = nextUrl.pathname;

  const authRoutes = ["/auth/login", "/auth/register"];

  if (authRoutes.some((path) => pathname.startsWith(path))) {
    return NextResponse.next();
  }

  const token = cookies.get("refresh_token")?.value;

  if (!token) {
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
