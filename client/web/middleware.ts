import { NextRequest, NextResponse } from "next/server";

export function middleware(req: NextRequest) {
  const { cookies, nextUrl } = req;
  const pathname = nextUrl.pathname;

  const token = cookies.get("refresh_token")?.value;

  const authRoutes = ["/auth/login", "/auth/register"];
  const protectedRoutes = ["/uptime"];
  const verifyRoutes = ["/auth/verify", "/auth/verify-success"];
  const resendRoute = "/auth/resend-verification";

  if (!token) {
    if (protectedRoutes.some((p) => pathname.startsWith(p))) {
      return NextResponse.redirect(new URL("/auth/login", req.url));
    }

    if (pathname.startsWith(resendRoute)) {
      return NextResponse.redirect(new URL("/auth/login", req.url));
    }

    return NextResponse.next();
  }

  let verified = false;
  try {
    const payloadBase64 = token.split(".")[1];
    const decodedPayload = JSON.parse(
      Buffer.from(payloadBase64, "base64").toString(),
    );
    verified = decodedPayload.verified === true;
  } catch {
    verified = false;
  }

  if (!verified && protectedRoutes.some((p) => pathname.startsWith(p))) {
    return NextResponse.redirect(new URL("/auth/resend-verification", req.url));
  }

  if (verified && authRoutes.some((p) => pathname.startsWith(p))) {
    return NextResponse.redirect(new URL("/uptime", req.url));
  }

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
