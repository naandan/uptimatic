"use client";

import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { urlService } from "@/lib/services/url";
import { formatDateGMT7, formatDateTimeGMT7, formatTimeGMT7 } from "@/utils/format";
import { ChevronLeft, ChevronRight, InfoIcon, Loader2 } from "lucide-react";
import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell } from "recharts";

type Bucket = {
  bucket_start: string;
  total_checks: number;
  up_checks: number;
  uptime_percent: number;
};

const getBarColor = (uptime: number) => {
  if (uptime >= 90) return "#22c55e";
  if (uptime >= 70) return "#f59e0b";
  return "#ef4444";
};

export default function UptimeStats() {
  const params = useParams();
  const id = params.id;
  const idNumber = Number(id);
  const router = useRouter();
  const [mode, setMode] = useState<"day" | "month">("day");
  const [offset, setOffset] = useState(0);
  const [data, setData] = useState<Bucket[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    const res = await urlService.stats(idNumber, mode, offset);
    if (res.success) {
      setData(res.data || []);
    } else {
      console.error(res.error);
      setData([]);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchData();
  }, [id, mode, offset]);

  const handlePrev = () => setOffset(offset + 1);
  const handleNext = () => setOffset(Math.max(0, offset - 1));

  return (
    <div className="mt-12 max-w-5xl mx-auto min-h-screen px-4">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <div className="flex items-center gap-2">
          <Button onClick={() => router.push('/uptime')} variant="ghost">
            <ChevronLeft className="w-4 h-4 mr-2" />
            {/* Kembali */}
          </Button>
          <h2 className="text-2xl font-semibold text-slate-800">Statistik Uptime</h2>
        </div>
        <div className="flex items-center gap-2">
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
      </div>

      {loading ? (
        <div className="flex flex-col items-center justify-center min-h-[50vh]">
          <Loader2 className="w-6 h-6 text-slate-500 animate-spin" />
          <p className="text-center text-slate-500 mt-2">Loading...</p>
        </div>
      ) : data.length === 0 ? (
        <div className="flex flex-col items-center justify-center min-h-[50vh]">
          <InfoIcon className="w-6 h-6 text-slate-500 mb-2" />
          <p className="text-center text-slate-500">Tidak ada data</p>
        </div>
      ) : (
        <ResponsiveContainer width="100%" height={400}>
          <BarChart data={data} margin={{ top: 20, right: 20, left: 0, bottom: 20 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="bucket_start"
              tickFormatter={(ts: string) => {
                return mode === "day" ? formatTimeGMT7(ts) : formatDateGMT7(ts);
              }}
            />
            <YAxis domain={[0, 100]} unit="%" />
            <Tooltip
              formatter={(value: any) => [`${value}%`, "Uptime"]}
              labelFormatter={(label: string) => `Waktu: ${mode === "day" ? formatDateTimeGMT7(label) : formatDateGMT7(label)}`}
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
