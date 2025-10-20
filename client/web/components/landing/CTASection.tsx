"use client";

import { Button } from "@/components/ui/button";
import Link from "next/link";

export const CTASection = () => {
  return (
    <section className="py-16 bg-gradient-to-r from-primary/90 to-primary text-white text-center">
      <div className="max-w-2xl mx-auto px-6">
        <h2 className="text-3xl font-semibold mb-4">
          Siap menjaga uptime Anda?
        </h2>
        <p className="text-white/90 mb-8">
          Mulai gratis sekarang dan nikmati pemantauan uptime otomatis selama 14
          hari.
        </p>
        <Button asChild size="lg" variant="secondary" className="font-semibold">
          <Link href="/auth/register">Mulai Sekarang</Link>
        </Button>
      </div>
    </section>
  );
};
