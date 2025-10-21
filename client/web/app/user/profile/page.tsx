"use client";

import { useState, useRef } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { toast } from "sonner";
import { useAuth } from "@/context/AuthContext";
import { userService } from "@/lib/services/user";
import { useRouter } from "next/navigation";
import Image from "next/image";

export default function ProfilePage() {
  const router = useRouter();
  const { user, setUser } = useAuth();
  const [name, setName] = useState(user?.name);
  const [email, setEmail] = useState(user?.email);
  const [open, setOpen] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);
  const [file, setFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement | null>(null);

  const handleFileSelect = (selectedFile: File) => {
    if (!selectedFile) return;

    if (selectedFile.size > 2 * 1024 * 1024) {
      toast.error("Ukuran file maksimal 2MB!");
      return;
    }

    if (!selectedFile.type.startsWith("image/")) {
      toast.error("Hanya file gambar yang diperbolehkan!");
      return;
    }

    setFile(selectedFile);
    setPreview(URL.createObjectURL(selectedFile));
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0];
    if (selectedFile) handleFileSelect(selectedFile);
  };

  const handleDragOver = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setIsDragging(false);
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) handleFileSelect(droppedFile);
  };

  const handleUpload = async () => {
    if (!file) return toast.error("Pilih gambar terlebih dahulu!");
    setIsUploading(true);

    try {
      const res = await userService.uploadURL(file.name, file.type);

      if (!res.success) throw new Error("Gagal meminta URL upload");

      const uploadUrl = res.data?.presigned_url;
      const fileName = res.data?.file_name;
      if (!uploadUrl || !fileName) throw new Error("URL upload tidak valid");

      const uploadRes = await fetch(uploadUrl, {
        method: "PUT",
        body: file,
      });

      if (!uploadRes.ok) throw new Error("Gagal upload ke storage");

      const updateRes = await userService.updateFoto(fileName);
      if (!updateRes.success) throw new Error("Gagal memperbarui foto profil");
      if (!updateRes.data) throw new Error("Gagal memperbarui foto profil");
      toast.success("Foto profil berhasil diperbarui!");
      if (updateRes.data?.url && user) {
        setPreview(updateRes.data?.url);
        setUser({ ...user, profile: updateRes.data?.url });
      }
      setOpen(false);
    } catch (err) {
      console.error(err);
      toast.error("Terjadi kesalahan saat upload gambar.");
    } finally {
      setIsUploading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name || !email) return toast.error("Semua field harus diisi!");
    const res = await userService.update(name, email);
    console.log(res.data);
    if (!res.success) {
      toast.error("Terjadi kesalahan saat memperbarui profil.");
      return;
    }
    if (user?.email != email) {
      router.push("/auth/resend-verification");
    }
    setUser(res.data);
    toast.success("Profil berhasil diperbarui!");
  };

  return (
    <div className="space-y-8">
      <h1 className="text-2xl font-bold">Profil Saya</h1>

      <div className="flex items-center gap-4">
        <Avatar
          className="w-16 h-16 cursor-pointer"
          onClick={() => setOpen(true)}
        >
          <AvatarImage
            src={
              user?.profile ||
              `https://placehold.co/100x100?text=${user?.name?.slice(0, 1)}`
            }
          />
          <AvatarFallback>ND</AvatarFallback>
        </Avatar>
        <div>
          <p className="font-semibold">{user?.name}</p>
          <p className="text-sm text-gray-500">{user?.email}</p>
        </div>
      </div>

      {/* Dialog Upload Gambar */}
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Ubah Foto Profil</DialogTitle>
          </DialogHeader>

          {/* Area Drag & Drop */}
          <div
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            onClick={() => fileInputRef.current?.click()}
            className={`border-2 border-dashed rounded-xl p-6 text-center cursor-pointer transition-colors
              ${isDragging ? "border-blue-500 bg-blue-50" : "border-gray-300 hover:border-gray-400"}
            `}
          >
            <input
              type="file"
              accept="image/*"
              ref={fileInputRef}
              onChange={handleFileChange}
              hidden
            />

            {preview ? (
              <div className="flex flex-col items-center space-y-3">
                <Image
                  width={100}
                  height={100}
                  src={preview}
                  alt="Preview"
                  className="w-28 h-28 rounded-full object-cover shadow"
                />
                <p className="text-sm text-gray-600">
                  Klik atau seret file lain untuk mengganti
                </p>
              </div>
            ) : (
              <div className="space-y-2">
                <p className="text-gray-700 font-medium">
                  {isDragging
                    ? "Lepas file di sini"
                    : "Seret atau klik untuk pilih foto"}
                </p>
                <p className="text-xs text-gray-500">
                  Maksimal ukuran 2 MB (JPG, PNG, WEBP)
                </p>
              </div>
            )}
          </div>

          <div className="mt-4 flex justify-end">
            <Button
              onClick={handleUpload}
              disabled={!file || isUploading}
              className="w-full md:w-auto"
            >
              {isUploading ? "Mengunggah..." : "Upload"}
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      {/* Form Profil */}
      <form className="space-y-4 max-w-lg" onSubmit={handleSubmit}>
        <div className="space-y-2">
          <Label htmlFor="name">Nama Lengkap</Label>
          <Input
            id="name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="email">Email</Label>
          <Input
            id="email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </div>

        <Button type="submit" className="w-full md:w-auto mt-1">
          Simpan Perubahan
        </Button>
      </form>
    </div>
  );
}
