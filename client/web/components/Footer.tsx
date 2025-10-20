"use client";

import Link from "next/link";
import { Github, Mail } from "lucide-react";

export const Footer = () => {
  return (
    <footer className="border-t border-slate-200 bg-white py-10 px-6">
      <div className="max-w-6xl mx-auto flex flex-col md:flex-row items-center justify-between gap-6 text-sm text-slate-600">
        <div className="text-center md:text-left">
          <h3 className="font-semibold text-slate-800 text-base mb-1">
            Uptimatic
          </h3>
          <p className="text-slate-500">
            Dibuat untuk keandalan web — pantau uptime Anda dengan presisi.
          </p>
        </div>

        <div className="flex items-center justify-center gap-6">
          <Link
            href="/privacy-policy"
            className="hover:text-primary transition"
          >
            Kebijakan Privasi
          </Link>
          <Link
            href="/terms-of-service"
            className="hover:text-primary transition"
          >
            Syarat & Ketentuan
          </Link>
        </div>

        <div className="flex items-center justify-center gap-4">
          <a
            href="https://github.com/naandan"
            target="_blank"
            rel="noopener noreferrer"
            className="hover:text-primary transition"
          >
            <Github className="w-5 h-5" />
          </a>
          <a
            href="mailto:nandanramdani608@gmail.com"
            className="hover:text-primary transition"
          >
            <Mail className="w-5 h-5" />
          </a>
        </div>
      </div>

      <div className="border-t border-slate-200 mt-8 pt-6 text-center text-slate-500 text-xs">
        © {new Date().getFullYear()} Uptimatic. Semua hak dilindungi.
      </div>
    </footer>
  );
};
