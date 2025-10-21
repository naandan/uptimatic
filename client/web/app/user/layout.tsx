"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Button } from "@/components/ui/button";
import { ChevronDown } from "lucide-react";

export default function UserLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  const menu = [
    { href: "/user/profile", label: "Akun Saya" },
    { href: "/user/reset-password", label: "Reset Kata Sandi" },
    { href: "/user/security", label: "Keamanan" },
  ];

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      {/* Mobile dropdown */}
      <div className="md:hidden mb-4">
        <Popover>
          <PopoverTrigger asChild>
            <Button variant="outline" className="w-[160px] justify-between">
              Menu Akun
              <ChevronDown className="ml-2 h-4 w-4" />
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto min-w-[160px] max-w-xs p-0">
            <nav className="flex flex-col text-sm">
              {menu.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`block px-4 py-2 hover:bg-gray-100 ${
                    pathname === item.href
                      ? "font-medium text-black"
                      : "text-gray-600"
                  }`}
                >
                  {item.label}
                </Link>
              ))}
            </nav>
          </PopoverContent>
        </Popover>
      </div>

      {/* Desktop sidebar + content */}
      <div className="hidden md:flex md:gap-8 min-h-[85vh]">
        {/* Sidebar */}
        <aside className="w-1/4 space-y-2">
          <h2 className="text-lg font-semibold mb-4">Navigasi Akun</h2>
          <nav className="flex flex-col space-y-2 text-sm">
            {menu.map((item) => (
              <Link
                key={item.href}
                href={item.href}
                className={`hover:underline ${
                  pathname === item.href
                    ? "font-medium text-black"
                    : "text-gray-600"
                }`}
              >
                {item.label}
              </Link>
            ))}
          </nav>
        </aside>

        {/* Main content */}
        <main className="flex-1">{children}</main>
      </div>

      {/* Mobile content fallback */}
      <div className="md:hidden">{children}</div>
    </div>
  );
}
