import { NextRequest, NextResponse } from "next/server";

export function middleware(req: NextRequest) {
  const { cookies, nextUrl } = req;
  const pathname = nextUrl.pathname;

  const authRoutes = ["/auth/login", "/auth/register"];

  // Route login/register → biarkan lewat
  if (authRoutes.some((path) => pathname.startsWith(path))) {
    return NextResponse.next();
  }

  // Ambil token dari cookie
  const token = cookies.get("refresh_token")?.value;

  if (!token) {
    return NextResponse.redirect(new URL("/auth/login", req.url));
  }

//   // Route admin → cek role
//   if (pathname.startsWith("/uptime")) {
//     try {
//       // Decode payload JWT tanpa memverifikasi signature
//       const payload = JSON.parse(Buffer.from(token.split(".")[1], "base64").toString());

//       if (payload.role !== "admin") {
//         // Role bukan admin → redirect ke home
//         return NextResponse.redirect(new URL("/", req.url));
//       }

//       // Role admin → lanjut
//       return NextResponse.next();
//     } catch (err) {
//       // Token invalid → redirect login
//       return NextResponse.redirect(new URL("/auth/login", req.url));
//     }
//   }

  // Semua route lain selain admin → token ada → lanjut
  return NextResponse.next();
}

// Tentukan route yang terkena middleware
export const config = {
  matcher: [
    // "/user/:path*",
    "/uptime/:path*",
    "/auth/login",
    "/auth/register",
  ],
};
