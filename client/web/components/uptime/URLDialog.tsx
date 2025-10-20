"use client";

import { Dispatch, SetStateAction } from "react";
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
import { URLResponse } from "@/types/url";
import { ErrorInput } from "@/types/response";
import ErrorInputMessage from "../ErrorInputMessage";

interface AddEditDialogProps {
  mode: "add" | "edit";
  open: boolean;
  initialData: Partial<URLResponse>;
  onInputChange: Dispatch<SetStateAction<Partial<URLResponse>>>;
  onSave: () => void;
  errors: ErrorInput[];
  onClose: () => void;
}

export function AddEditDialog({
  mode,
  open,
  initialData,
  onInputChange,
  onSave,
  errors,
  onClose,
}: AddEditDialogProps) {
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
              value={initialData.label}
              onChange={(e) =>
                onInputChange({ ...initialData, label: e.target.value })
              }
            />
            <ErrorInputMessage errors={errors} field="label" />
          </div>

          <div className="space-y-2">
            <Label htmlFor="url">URL</Label>
            <Input
              id="url"
              type="url"
              placeholder="https://www.example.com"
              value={initialData.url}
              onChange={(e) =>
                onInputChange({ ...initialData, url: e.target.value })
              }
            />
            <ErrorInputMessage errors={errors} field="url" />
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
              checked={initialData.active}
              onCheckedChange={(checked) =>
                onInputChange({ ...initialData, active: checked ?? false })
              }
            />
          </div>
          <ErrorInputMessage errors={errors} field="active" />
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Batal
          </Button>
          <Button onClick={onSave}>
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
