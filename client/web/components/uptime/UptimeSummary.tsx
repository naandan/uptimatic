"use client";

import { Check, X } from "lucide-react";

export function UptimeSummary({
  avgUptime,
  totalChecks,
  upChecks,
  downChecks,
  dateRange,
}: {
  avgUptime: number;
  totalChecks: number;
  upChecks: number;
  downChecks: number;
  dateRange: string;
}) {
  const summary = [
    { title: "Periode", value: dateRange, color: "text-black" },
    {
      title: "Rata-rata Uptime",
      value: `${avgUptime}%`,
      color: "text-green-600",
    },
    { title: "Total Pemeriksaan", value: totalChecks, color: "text-slate-800" },
    {
      title: "UP / DOWN",
      value: (
        <>
          <p className="text-green-600 flex items-center">
            <Check className="inline-block w-6 h-6 mr-1" />
            {upChecks}
          </p>
          <p className="text-red-500 flex items-center">
            <X className="inline-block w-6 h-6 mr-1" />
            {downChecks}
          </p>
        </>
      ),
    },
  ];

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8 lg:border-b-2 pb-6">
      {summary.map((item, i) => (
        <div
          key={i}
          className="flex flex-col items-center justify-center gap-1 border-b-2 pb-6 lg:pb-0 lg:border-b-0 lg:border-r-2 lg:last:border-none"
        >
          <div className="text-slate-600 mb-1">{item.title}</div>
          <div className={`text-lg font-semibold ${item.color || ""}`}>
            {item.value}
          </div>
        </div>
      ))}
    </div>
  );
}
