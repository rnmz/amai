import { api, BASE_URL } from "./client";
import type { ApiMessage } from "./client";

export interface FileItem {
  id: string;
  ext: string;
}

export interface FileListResponse {
  files: FileItem[];
  pages: number;
}


export function getFileUrl(id: string): string {
  return `${BASE_URL}/file/get?id=${encodeURIComponent(id)}`;
}

export async function getFileBlob(id: string): Promise<Blob> {
  const { data } = await api.get<Blob>("/file/get", {
    params: { id },
    responseType: "blob",
  });
  return data;
}

export async function getFilesList(page: number): Promise<FileListResponse> {
  const { data } = await api.get<FileListResponse>("/admin/file/list", {
    params: { page },
  });
  return data;
}

export interface UploadFileResponse extends ApiMessage {
  file_id: string;
}

export async function uploadFile(
  file: File,
  onProgress?: (percent: number) => void
): Promise<UploadFileResponse> {
  const formData = new FormData();
  formData.append("file", file);

  const { data } = await api.post<UploadFileResponse>("/admin/file/upload", formData, {
    onUploadProgress: (event) => {
      if (onProgress && event.total) {
        onProgress(Math.round((event.loaded * 100) / event.total));
      }
    },
  });
  return data;
}

export async function deleteFile(id: string): Promise<ApiMessage> {
  const { data } = await api.delete<ApiMessage>("/admin/file/delete", {
    params: { id },
  });
  return data;
}