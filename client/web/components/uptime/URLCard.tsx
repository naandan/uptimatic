"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { URL } from "@/types/uptime";
import { Globe, Edit, Trash2 } from "lucide-react";

interface URLCardProps {
  data: URL;
  onToggle: (id: number) => void;
  onEdit: () => void;
  onDelete: () => void;
}


export default function URLCard({ data, onToggle, onEdit, onDelete }: URLCardProps) {
  return (
    <div className="p-5 bg-white rounded-2xl border border-slate-200 shadow-sm hover:shadow-md transition flex flex-col justify-between">
      <div>
        <div className="flex items-center justify-between mb-3">
          <div className="flex items-center gap-2">
            <Globe className="w-5 h-5 text-primary" />
            <h3 className="font-semibold text-slate-800">{data.label}</h3>
          </div>
          <Badge
            variant={data.active ? "default" : "secondary"}
            className={
              data.active ? "bg-green-500 hover:bg-green-600" : "bg-gray-400"
            }
          >
            {data.active ? "Aktif" : "Nonaktif"}
          </Badge>
        </div>

        <p className="text-sm text-primary truncate mb-2">{data.url}</p>
        <p className="text-xs text-slate-500">
          Interval: {data.interval} menit
        </p>
        <p className="text-xs text-slate-500">
          Terakhir dicek:{" "}
          {new Date(data.last_checked).toLocaleString("id-ID", {
            hour: "2-digit",
            minute: "2-digit",
          })}
        </p>
      </div>

      <div className="flex items-center justify-between mt-4">
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={onEdit}>
            <Edit className="w-4 h-4" />
          </Button>
          <Button variant="outline" size="sm" onClick={onDelete}>
            <Trash2 className="w-4 h-4 text-red-500" />
          </Button>
        </div>
        <Button
          variant={data.active ? "secondary" : "default"}
          size="sm"
          onClick={() => onToggle(data.id)}
        >
          {data.active ? "Nonaktifkan" : "Aktifkan"}
        </Button>
      </div>
    </div>
  );
}
