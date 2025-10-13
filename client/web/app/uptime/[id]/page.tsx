"use client";

import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ChevronLeft, ChevronRight } from "lucide-react";
import { useState, useEffect } from "react";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from "recharts";

type Bucket = {
  bucket_start: string;
  total_checks: number;
  up_checks: number;
  uptime_percent: number;
};

interface UptimeData {
  data: Bucket[];
  status: string;
  request_id: string;
}

interface Props {
  urlId: number; // ID URL yang akan ditampilkan
}

const getBarColor = (uptime: number) => {
  if (uptime >= 90) return "#22c55e"; // hijau
  if (uptime >= 70) return "#f59e0b"; // kuning/orange
  return "#ef4444"; // merah
};

// utils/generateDummyUptime.ts
export function generateDummyUptime(mode: "day" | "month", offset: number = 0) {
  const data = [];
  const now = new Date();

  if (mode === "day") {
    // per jam
    const day = new Date();
    day.setDate(now.getDate() - offset);

    for (let h = 0; h < 24; h++) {
      const total_checks = Math.floor(Math.random() * 20) + 1;
      const up_checks = Math.floor(Math.random() * (total_checks + 1));
      data.push({
        bucket_start: new Date(day.getFullYear(), day.getMonth(), day.getDate(), h).toISOString(),
        total_checks,
        up_checks,
        uptime_percent: total_checks === 0 ? 100 : parseFloat(((up_checks / total_checks) * 100).toFixed(2)),
      });
    }
  } else if (mode === "month") {
    // per hari
    const year = now.getFullYear();
    const month = now.getMonth() - offset;
    const daysInMonth = new Date(year, month + 1, 0).getDate();

    for (let d = 1; d <= daysInMonth; d++) {
      const total_checks = Math.floor(Math.random() * 100) + 1;
      const up_checks = Math.floor(Math.random() * (total_checks + 1));
      data.push({
        bucket_start: new Date(year, month, d).toISOString(),
        total_checks,
        up_checks,
        uptime_percent: parseFloat(((up_checks / total_checks) * 100).toFixed(2)),
      });
    }
  }

  return { data, status: "success", request_id: crypto.randomUUID() };
}


export default function UptimeStats({ urlId }: Props) {
  const [mode, setMode] = useState<"day" | "month">("day");
  const [offset, setOffset] = useState(0);
  const [data, setData] = useState<Bucket[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
    //   const res = await fetch(`/api/uptime?id=${urlId}&mode=${mode}&offset=${offset}`);
    //   const json: UptimeData = await res.json();
      setData(generateDummyUptime("day", offset).data);
    } catch (err) {
      console.error(err);
      setData([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
  }, [urlId, mode, offset]);

  const handlePrev = () => setOffset(offset + 1);
  const handleNext = () => setOffset(Math.max(0, offset - 1));

  return (
    <div className="p-6 max-w-5xl mx-auto">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <div className="flex gap-2">
          <Button onClick={handlePrev} variant="outline">
            <ChevronLeft className="w-4 h-4" />
          </Button>
          <Button onClick={handleNext} variant="outline">
            <ChevronRight className="w-4 h-4" />
          </Button>
        </div>

        <Select value={mode} onValueChange={(v: any) => setMode(v)}>
          <SelectTrigger>
            <SelectValue placeholder="Mode" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="day">Per Hari</SelectItem>
            <SelectItem value="month">Per Bulan</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {loading ? (
        <p className="text-center text-slate-500">Loading...</p>
      ) : data.length === 0 ? (
        <p className="text-center text-slate-500">Tidak ada data</p>
      ) : (
        <ResponsiveContainer width="100%" height={400}>
          <BarChart data={data} margin={{ top: 20, right: 20, left: 0, bottom: 20 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="bucket_start"
              tickFormatter={(ts: string) => {
                const d = new Date(ts);
                return mode === "day"
                  ? `${d.getHours()}:00`
                  : `${d.getDate()}/${d.getMonth() + 1}`;
              }}
            />
            <YAxis domain={[0, 100]} unit="%" />
            <Tooltip
              formatter={(value: any) => [`${value}%`, "Uptime"]}
              labelFormatter={(label: string) => `Waktu: ${new Date(label).toLocaleString()}`}
            />
            <Bar dataKey="uptime_percent" radius={[4, 4, 0, 0]}>
              {data.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={getBarColor(entry.uptime_percent)} />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      )}
    </div>
  );
}
