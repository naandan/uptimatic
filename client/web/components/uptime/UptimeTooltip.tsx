"use client";

import { FC } from "react";

interface Props {
  active?: boolean;
  payload?: { value: number }[];
  label?: string | number;
  mode: "day" | "month" | "year";
}

export const UptimeTooltip: FC<Props> = ({ active, payload, label, mode }) => {
  if (!active || !payload?.length) return null;

  const uptime = payload[0].value;
  const date = new Date(label ?? "");

  let formattedLabel = "";
  if (mode === "day") {
    formattedLabel = date.toLocaleString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
      timeZone: "Asia/Jakarta",
    });
  } else if (mode === "month") {
    formattedLabel = date.toLocaleDateString("id-ID", {
      day: "2-digit",
      month: "long",
      timeZone: "Asia/Jakarta",
    });
  } else {
    formattedLabel = date.toLocaleDateString("id-ID", {
      month: "long",
      year: "numeric",
      timeZone: "Asia/Jakarta",
    });
  }

  return (
    <div className="rounded-lg border bg-white p-3 shadow-sm text-slate-700">
      <div className="text-xs text-slate-500 mb-1">{formattedLabel}</div>
      <div className="flex items-center gap-2">
        <div
          className="h-2.5 w-2.5 rounded-full"
          style={{
            backgroundColor:
              uptime >= 90 ? "#22c55e" : uptime >= 70 ? "#f59e0b" : "#ef4444",
          }}
        />
        <span className="text-sm font-medium">Uptime:</span>
        <span className="text-sm font-semibold">{uptime.toFixed(2)}%</span>
      </div>
    </div>
  );
};
