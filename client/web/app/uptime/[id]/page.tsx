"use client";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { urlService } from "@/lib/services/url";
import { URLResponse, URLStats } from "@/types/url";
import {
  formatDateGMT7,
  formatDateTimeGMT7,
  formatMonthGMT7,
  formatTimeGMT7,
} from "@/utils/format";
import {
  ChevronLeft,
  ChevronRight,
  Globe,
  Link2,
  CheckCircle2,
  XCircle,
  Loader2,
  InfoIcon,
  Check,
  X,
  RotateCw,
} from "lucide-react";
import { useParams, useRouter } from "next/navigation";
import { useState, useEffect, useCallback } from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import { toast } from "sonner";
import { motion } from "framer-motion";

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
  const [data, setData] = useState<URLStats[]>([]);
  const [loading, setLoading] = useState(false);
  const [url, setUrl] = useState<URLResponse>();

  const fetchURL = useCallback(async () => {
    const res = await urlService.get(idNumber);
    if (!res.success) {
      toast.error("URL tidak ditemukan");
      return;
    }
    if (!res.data) return;
    setUrl(res.data);
  }, [idNumber]);

  const fetchStats = useCallback(async () => {
    setLoading(true);
    const res = await urlService.stats(idNumber, mode, offset);
    if (res.success) {
      setData(res.data || []);
    } else {
      console.error(res.error);
      setData([]);
    }
    setLoading(false);
  }, [idNumber, mode, offset]);

  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  useEffect(() => {
    fetchURL();
  }, [fetchURL]);

  const fetchAll = async () => {
    await fetchStats();
    await fetchURL();
  };

  const handlePrev = () => setOffset(offset + 1);
  const handleNext = () => setOffset(Math.max(0, offset - 1));

  const avgUptime =
    data.length > 0
      ? (data.reduce((a, b) => a + b.uptime_percent, 0) / data.length).toFixed(
          2,
        )
      : 0;

  const totalChecks = data.reduce((a, b) => a + b.total_checks, 0);
  const upChecks = data.reduce((a, b) => a + b.up_checks, 0);
  const downChecks = totalChecks - upChecks;

  const dateRange =
    data.length > 0
      ? mode === "day"
        ? `${formatDateGMT7(data[0].bucket_start)}`
        : `${formatMonthGMT7(data[0].bucket_start)}`
      : "-";

  return (
    <div className="mt-12 max-w-5xl mx-auto min-h-screen px-4 mb-14">
      <motion.div
        className="p-6 mb-8 border-b-2"
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
      >
        <div className="flex items-center justify-between flex-wrap gap-4">
          <div className="flex gap-3 flex-col items-start sm:flex-row sm:items-center">
            <Button onClick={() => router.push("/uptime")} variant="outline">
              <ChevronLeft className="w-5 h-5 mr-1" />
            </Button>
            <div>
              <h2 className="text-2xl font-bold text-slate-800">
                Statistik Uptime
              </h2>
              {url ? (
                <div>
                  <div className="flex items-center gap-2 mt-2 text-slate-600 text-sm">
                    <Link2 className="w-4 h-4 text-slate-500" />
                    <a
                      href={url.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="hover:underline"
                    >
                      {url.label} — {url.url}
                    </a>
                    {url.active ? (
                      <span className="flex items-center text-green-600 text-xs ml-2">
                        <CheckCircle2 className="w-3 h-3 mr-1" /> Aktif
                      </span>
                    ) : (
                      <span className="flex items-center text-red-500 text-xs ml-2">
                        <XCircle className="w-3 h-3 mr-1" /> Nonaktif
                      </span>
                    )}
                  </div>

                  <div className="mt-2 text-xs text-slate-500 flex flex-wrap items-center gap-4">
                    <p>
                      <Globe className="inline-block w-3 h-3 mr-1" /> Terakhir
                      diperiksa:{" "}
                      <span className="font-medium text-slate-700">
                        {formatDateTimeGMT7(url.last_checked)}
                      </span>
                    </p>
                  </div>
                </div>
              ) : (
                <p className="text-sm text-slate-400">Memuat info URL...</p>
              )}
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button onClick={fetchAll} variant="outline" size="icon">
              <RotateCw className="w-6 h-6" />
            </Button>
            <Button onClick={handlePrev} variant="outline" size="icon">
              <ChevronLeft className="w-6 h-6" />
            </Button>
            <Button onClick={handleNext} variant="outline" size="icon">
              <ChevronRight className="w-6 h-6" />
            </Button>
            <Select
              value={mode}
              onValueChange={(v: "day" | "month") => {
                setMode(v);
                setOffset(0);
              }}
            >
              <SelectTrigger className="w-32">
                <SelectValue placeholder="Mode" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="day">Per Hari</SelectItem>
                <SelectItem value="month">Per Bulan</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </motion.div>

      {!loading && data.length > 0 && (
        <motion.div
          className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8 lg:border-b-2 pb-6"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
        >
          {[
            { title: "Periode", value: dateRange, color: "text-black" },
            {
              title: "Rata-rata Uptime",
              value: `${avgUptime}%`,
              color: "text-green-600",
            },
            {
              title: "Total Pemeriksaan",
              value: totalChecks,
              color: "text-slate-800",
            },
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
          ].map((item, i) => (
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
        </motion.div>
      )}

      {/* Chart Section */}
      {loading ? (
        <div className="flex flex-col items-center justify-center min-h-[50vh]">
          <Loader2 className="w-6 h-6 text-slate-500 animate-spin" />
          <p className="text-center text-slate-500 mt-2">Memuat data...</p>
        </div>
      ) : data.length === 0 ? (
        <div className="flex flex-col items-center justify-center min-h-[50vh]">
          <InfoIcon className="w-6 h-6 text-slate-500 mb-2" />
          <p className="text-center text-slate-500">Tidak ada data</p>
        </div>
      ) : (
        <motion.div
          className="bg-white rounded-2xl p-3 border-slate-100"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
        >
          <h3 className="text-lg font-semibold text-slate-700 mb-4">
            Grafik Uptime
          </h3>
          <ResponsiveContainer width="100%" height={400}>
            <BarChart
              data={data}
              margin={{ top: 20, right: 20, left: 0, bottom: 20 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="bucket_start"
                tickFormatter={(ts: string) =>
                  mode === "day" ? formatTimeGMT7(ts) : formatDateGMT7(ts)
                }
              />
              <YAxis domain={[0, 100]} unit="%" />
              <Tooltip
                contentStyle={{ borderRadius: "8px" }}
                formatter={(value: number) => [`${value}%`, "Uptime"]}
                labelFormatter={(label: string) =>
                  mode === "day"
                    ? formatDateTimeGMT7(label)
                    : formatDateGMT7(label)
                }
              />
              <Bar dataKey="uptime_percent" radius={[4, 4, 0, 0]}>
                {data.map((entry, index) => (
                  <Cell
                    key={`cell-${index}`}
                    fill={getBarColor(entry.uptime_percent)}
                  />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>

          {/* Legend */}
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
        </motion.div>
      )}
    </div>
  );
}
