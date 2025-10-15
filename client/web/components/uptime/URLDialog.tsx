"use client";

import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { URL } from "@/types/url";

interface AddEditDialogProps {
  mode: "add" | "edit";
  open: boolean;
  initialData?: Partial<URL>;
  onClose: () => void;
  onSave: (data: Omit<URL, "id" | "last_checked" | "created_at">) => void;
}

export function AddEditDialog({
  mode,
  open,
  initialData,
  onClose,
  onSave,
}: AddEditDialogProps) {
  const [form, setForm] = useState({
    label: initialData?.label || "",
    url: initialData?.url || "",
    interval: initialData?.interval || 5,
    active: initialData?.active ?? true,
  });

  const handleSubmit = () => {
    if (!form.label.trim() || !form.url.trim()) return;
    onSave({
      label: form.label.trim(),
      url: form.url.trim(),
      interval: Number(form.interval),
      active: form.active,
    });

    setForm({
      label: "",
      url: "",
      interval: 5,
      active: true,
    });
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>
            {mode === "add" ? "Tambah URL Baru" : "Edit URL"}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-2">
          <div className="space-y-2">
            <Label htmlFor="label">Label</Label>
            <Input
              id="label"
              placeholder="Contoh: Homepage"
              value={form.label}
              onChange={(e) => setForm({ ...form, label: e.target.value })}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="url">URL</Label>
            <Input
              id="url"
              type="url"
              placeholder="https://www.example.com"
              value={form.url}
              onChange={(e) => setForm({ ...form, url: e.target.value })}
            />
          </div>

          {/* <div className="space-y-2">
            <Label htmlFor="interval">Interval (menit)</Label>
            <Input
              id="interval"
              type="number"
              min={1}
              value={form.interval}
              onChange={(e) =>
                setForm({ ...form, interval: Number(e.target.value) })
              }
            />
          </div> */}

          <div className="flex items-center justify-between">
            <Label>Aktif</Label>
            <Switch
              checked={form.active}
              onCheckedChange={(checked) =>
                setForm({ ...form, active: checked?? false })
              }
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Batal
          </Button>
          <Button onClick={handleSubmit}>
            {mode === "add" ? "Tambah" : "Simpan"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

interface DeleteDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  targetLabel?: string;
}

export function DeleteDialog({
  open,
  onClose,
  onConfirm,
  targetLabel,
}: DeleteDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Hapus URL</DialogTitle>
        </DialogHeader>
        <p className="text-sm text-slate-600">
          Apakah Anda yakin ingin menghapus{" "}
          <span className="font-semibold text-slate-800">{targetLabel}</span>?
          Tindakan ini tidak dapat dibatalkan.
        </p>
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Batal
          </Button>
          <Button variant="destructive" onClick={onConfirm}>
            Hapus
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
