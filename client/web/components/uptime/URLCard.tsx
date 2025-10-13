"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { URL } from "@/types/uptime";
import { formatTimeGMT7 } from "@/utils/format";
import { Globe, Trash2, Edit2 } from "lucide-react";
import { useRouter } from "next/navigation";

interface URLCardProps {
  data: URL;
  onToggle: (id: number) => void;
  onEdit: () => void;
  onDelete: () => void;
}

export default function URLCard({ data, onToggle, onEdit, onDelete }: URLCardProps) {
  const router = useRouter();

  const handleCardClick = () => {
    router.push(`/uptime/${data.id}`);
  };

  const stopClick = (e: React.MouseEvent) => e.stopPropagation();

  return (
    <div
      onClick={handleCardClick}
      className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition flex flex-col justify-between cursor-pointer"
    >
      <div>
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <Globe className="w-5 h-5 text-primary" />
            <h3 className="font-semibold text-slate-800">{data.label}</h3>
          </div>
          <Badge
            variant={data.active ? "default" : "secondary"}
            className={data.active ? "bg-green-500 hover:bg-green-600" : "bg-gray-400"}
          >
            {data.active ? "Aktif" : "Nonaktif"}
          </Badge>
        </div>

        <p className="text-sm text-primary truncate mb-2">{data.url}</p>
        <p className="text-xs text-slate-500">Interval: {data.interval / 60} menit</p>
        <p className="text-xs text-slate-500">
          Terakhir dicek: {formatTimeGMT7(data.last_checked)}
        </p>
      </div>

      <div className="flex items-center justify-between mt-4">
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={(e) => { stopClick(e); onEdit(); }}>
            <Edit2 className="w-4 h-4" />
          </Button>
          <Button variant="outline" size="sm" onClick={(e) => { stopClick(e); onDelete(); }}>
            <Trash2 className="w-4 h-4 text-red-500" />
          </Button>
        </div>
        <Button
          variant={data.active ? "secondary" : "default"}
          size="sm"
          onClick={(e) => {
            stopClick(e);
            onToggle(data.id);
          }}
        >
          {data.active ? "Nonaktifkan" : "Aktifkan"}
        </Button>
      </div>
    </div>
  );
}
