"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
} from "recharts";
import {
  formatHourShort,
  formatDayShort,
  formatMonthShort,
} from "@/utils/format";
import { UptimeTooltip } from "./UptimeTooltip";
import { URLStats } from "@/types/url";

interface Props {
  data: URLStats[];
  mode: "day" | "month" | "year";
}

export const UptimeChart = ({ data, mode }: Props) => {
  return (
    <div className="bg-white rounded-2xl p-3 border border-slate-100">
      <h3 className="text-lg font-semibold text-slate-700 mb-2">
        Grafik Uptime
      </h3>

      <ResponsiveContainer width="100%" height={400}>
        <AreaChart
          data={data}
          margin={{ top: 20, right: 20, left: 0, bottom: 20 }}
        >
          <defs>
            <linearGradient id="fillUptime" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#22c55e" stopOpacity={0.8} />
              <stop offset="95%" stopColor="#22c55e" stopOpacity={0.1} />
            </linearGradient>
          </defs>

          <CartesianGrid strokeDasharray="3 3" vertical={false} />

          <XAxis
            dataKey="bucket_start"
            tickLine={false}
            axisLine={false}
            tickMargin={8}
            minTickGap={32}
            tickFormatter={(value) => {
              if (mode === "day") return formatHourShort(value);
              if (mode === "month") return formatDayShort(value);
              return formatMonthShort(value);
            }}
          />

          <YAxis domain={[0, 100]} tickFormatter={(v) => `${v}%`} />

          <Tooltip content={<UptimeTooltip mode={mode} />} />

          <Area
            type="monotone"
            dataKey="uptime_percent"
            stroke="#22c55e"
            fill="url(#fillUptime)"
            strokeWidth={2}
            dot={false}
          />
        </AreaChart>
      </ResponsiveContainer>

      <div className="flex justify-center gap-4 mt-4 text-sm text-slate-600">
        <div className="flex items-center gap-1">
          <div className="w-4 h-4 bg-[#22c55e] rounded-sm" /> &ge; 90%
        </div>
        <div className="flex items-center gap-1">
          <div className="w-4 h-4 bg-[#f59e0b] rounded-sm" /> 70–89%
        </div>
        <div className="flex items-center gap-1">
          <div className="w-4 h-4 bg-[#ef4444] rounded-sm" /> &lt; 70%
        </div>
      </div>
    </div>
  );
};
