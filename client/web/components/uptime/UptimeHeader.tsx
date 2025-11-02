"use client";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import {
  ChevronLeft,
  Link2,
  CheckCircle2,
  XCircle,
  RotateCw,
  CalendarIcon,
  Globe,
} from "lucide-react";
import { formatDateGMT7, formatDateTimeGMT7 } from "@/utils/format";
import { URLResponse } from "@/types/url";

export function UptimeHeader({
  url,
  mode,
  date,
  onDateChange,
  onModeChange,
  onRefresh,
  onBack,
}: {
  url?: URLResponse;
  mode: "day" | "month" | "year";
  date: Date;
  onDateChange: (d: Date | undefined) => void;
  onModeChange: (v: "day" | "month" | "year") => void;
  onRefresh: () => void;
  onBack: () => void;
}) {
  return (
    <div className="p-6 mb-8 border-b-2 flex flex-wrap items-center justify-between gap-4">
      <div className="flex gap-3 flex-col items-start sm:flex-row sm:items-center">
        <Button onClick={onBack} variant="outline">
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

      <div className="flex flex-wrap items-center gap-2">
        <Button onClick={onRefresh} variant="outline" size="icon">
          <RotateCw className="w-6 h-6" />
        </Button>

        <Popover>
          <PopoverTrigger asChild>
            <Button
              variant="outline"
              className="flex-1 min-w-[180px] sm:w-[240px] justify-start text-left font-normal"
            >
              <CalendarIcon className="mr-2 h-4 w-4" />
              {date ? (
                formatDateGMT7(date.toString())
              ) : (
                <span>Pilih tanggal</span>
              )}
            </Button>
          </PopoverTrigger>
          <PopoverContent className="w-auto p-0">
            <Calendar
              mode="single"
              selected={date}
              onSelect={onDateChange}
              captionLayout="dropdown"
              className="rounded-md border shadow-sm"
              required
            />
          </PopoverContent>
        </Popover>

        <Select
          value={mode}
          onValueChange={(v: "day" | "month" | "year") => onModeChange(v)}
        >
          <SelectTrigger className="w-full sm:w-32">
            <SelectValue placeholder="Mode" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="day">Hari ini</SelectItem>
            <SelectItem value="month">Bulan ini</SelectItem>
            <SelectItem value="year">Tahun ini</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
