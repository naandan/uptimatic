"use client";

import { CheckCircle2, BarChart3, Bell, Clock } from "lucide-react";

const features = [
  {
    icon: <Clock className="w-6 h-6 text-primary" />,
    title: "Monitoring Otomatis",
    desc: "Pantau uptime situs atau API Anda setiap menit dengan akurasi tinggi.",
  },
  {
    icon: <Bell className="w-6 h-6 text-primary" />,
    title: "Notifikasi Real-Time",
    desc: "Dapatkan notifikasi email atau Slack ketika downtime terdeteksi.",
  },
  {
    icon: <BarChart3 className="w-6 h-6 text-primary" />,
    title: "Statistik Performa",
    desc: "Analisis tren uptime dan respon server secara mendetail.",
  },
  {
    icon: <CheckCircle2 className="w-6 h-6 text-primary" />,
    title: "Keandalan 99.9%",
    desc: "Infrastruktur global untuk pemantauan yang cepat dan stabil.",
  },
];

export const FeaturesSection = () => {
  return (
    <section id="features" className="py-24 bg-white">
      <div className="max-w-5xl mx-auto px-6 text-center">
        <h2 className="text-3xl font-semibold text-slate-800 mb-4">
          Fitur Utama Uptimatic
        </h2>
        <p className="text-slate-600 mb-12">
          Semua alat yang Anda butuhkan untuk menjaga layanan Anda tetap online.
        </p>
        <div className="grid md:grid-cols-2 gap-3 md:gap-6">
          {features.map((f, i) => (
            <div
              key={i}
              className="flex flex-col items-center text-center space-y-3 p-6 rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition"
            >
              {f.icon}
              <h3 className="font-semibold text-lg">{f.title}</h3>
              <p className="text-slate-600 text-sm">{f.desc}</p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
