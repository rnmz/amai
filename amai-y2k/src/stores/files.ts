import { defineStore } from "pinia";
import { ref } from "vue";
import {
  getFilesList,
  uploadFile,
  deleteFile,
  getFileUrl,
  type FileItem,
} from "@/api/file";
import { getApiErrorMessage } from "@/api/client";

export const useFilesStore = defineStore("files", () => {
  const files = ref<FileItem[]>([]);
  const pages = ref(0);
  const currentPage = ref(1);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  const isUploading = ref(false);
  const uploadProgress = ref(0);

  async function fetchFiles(page = currentPage.value): Promise<void> {
    isLoading.value = true;
    error.value = null;
    try {
      const data = await getFilesList(page);
      files.value = data.files;
      pages.value = data.pages;
      currentPage.value = page;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to load file list");
    } finally {
      isLoading.value = false;
    }
  }

  async function upload(file: File): Promise<string | null> {
    isUploading.value = true;
    uploadProgress.value = 0;
    error.value = null;
    try {
      const result = await uploadFile(file, (percent) => {
        uploadProgress.value = percent;
      });
      return result.file_id;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to upload file");
      return null;
    } finally {
      isUploading.value = false;
    }
  }

  async function remove(id: string): Promise<boolean> {
    isLoading.value = true;
    error.value = null;
    try {
      await deleteFile(id);
      files.value = files.value.filter((f) => f.id !== id);
      return true;
    } catch (e) {
      error.value = getApiErrorMessage(e, "Failed to delete file");
      return false;
    } finally {
      isLoading.value = false;
    }
  }

  return {
    files,
    pages,
    currentPage,
    isLoading,
    error,
    isUploading,
    uploadProgress,
    fetchFiles,
    upload,
    remove,
    getFileUrl,
  };
});