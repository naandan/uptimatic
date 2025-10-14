"use client";

import { Button } from "@/components/ui/button";
import Link from "next/link";
import { motion } from "framer-motion";

export const HeroSection = () => {
  return (
    <section className="pt-36 pb-28 text-center bg-gradient-to-b from-slate-50 to-white">
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.6 }}
        className="max-w-3xl mx-auto px-4"
      >
        <h1 className="text-5xl font-bold leading-tight tracking-tight text-slate-800">
          Pantau Ketersediaan Website Anda.  
          <span className="text-primary"> Secara Real-Time.</span>
        </h1>
        <p className="mt-6 text-lg text-slate-600">
          Uptimatic membantu Anda memonitor uptime dan performa situs atau API,  
          agar Anda selalu tahu sebelum pengguna Anda tahu.
        </p>
        <div className="mt-8 flex justify-center gap-4">
          <Button asChild size="lg">
            <Link href="/auth/register">Mulai Gratis</Link>
          </Button>
          <Button variant="outline" asChild size="lg">
            <Link href="#features">Lihat Fitur</Link>
          </Button>
        </div>
      </motion.div>
    </section>
  );
}
