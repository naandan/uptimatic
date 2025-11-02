"use client";
import { useEffect, useState, useCallback } from "react";
import { toast } from "sonner";
import { useParams, useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { Loader2, InfoIcon } from "lucide-react";
import { urlService } from "@/lib/services/url";
import {
  formatDateGMT7,
  formatMonthGMT7,
  formatYearGMT7,
} from "@/utils/format";
import { UptimeHeader } from "@/components/uptime/UptimeHeader";
import { UptimeSummary } from "@/components/uptime/UptimeSummary";
import { UptimeChart } from "@/components/uptime/UptimeChart";
import { URLResponse, URLStats } from "@/types/url";

export default function UptimeContainer() {
  const { id } = useParams();
  const router = useRouter();
  const publicId = String(id);
  const [mode, setMode] = useState<"day" | "month" | "year">("day");
  const [date, setDate] = useState<Date>(new Date());
  const [data, setData] = useState<URLStats[]>([]);
  const [loading, setLoading] = useState(false);
  const [url, setUrl] = useState<URLResponse | null>();

  const fetchURL = useCallback(async () => {
    const res = await urlService.get(publicId);
    if (!res.success) toast.error("URL tidak ditemukan");
    else setUrl(res.data);
  }, [publicId]);

  const fetchStats = useCallback(async () => {
    setLoading(true);
    const res = await urlService.stats(publicId, mode, date);
    setData(res.success ? res.data || [] : []);
    setLoading(false);
  }, [publicId, mode, date]);

  useEffect(() => {
    fetchURL();
  }, [fetchURL]);
  useEffect(() => {
    fetchStats();
  }, [fetchStats]);

  const totalChecks = data.reduce((a, b) => a + b.total_checks, 0);
  const upChecks = data.reduce((a, b) => a + b.up_checks, 0);
  const downChecks = totalChecks - upChecks;
  const avgUptime =
    totalChecks > 0 ? Math.round((upChecks * 10000) / totalChecks) / 100 : 0;

  const dateRange = data.length
    ? mode === "day"
      ? formatDateGMT7(data[0].bucket_start)
      : mode === "month"
        ? formatMonthGMT7(data[0].bucket_start)
        : formatYearGMT7(data[0].bucket_start)
    : "-";

  return (
    <div className="mt-12 max-w-5xl mx-auto min-h-screen px-4 mb-14">
      <UptimeHeader
        url={url || undefined}
        mode={mode}
        date={date}
        onModeChange={setMode}
        onDateChange={(d) => d && setDate(d)}
        onRefresh={() => {
          fetchStats();
          fetchURL();
        }}
        onBack={() => router.push("/uptime")}
      />

      {!loading && data.length > 0 && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <UptimeSummary
            avgUptime={avgUptime}
            totalChecks={totalChecks}
            upChecks={upChecks}
            downChecks={downChecks}
            dateRange={dateRange}
          />
        </motion.div>
      )}

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
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <UptimeChart data={data} mode={mode} />
        </motion.div>
      )}
    </div>
  );
}
