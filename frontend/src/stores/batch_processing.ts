import { BatchData } from "@/api/batch";
import { create } from "zustand";
import { FileStoreType } from "@/stores/file_store";

export const useBatchFilesStore = create<FileStoreType>((set) => ({
  files: [],
  setFiles: (files) => set({ files: files }),
  clearFiles: () => set({ files: [] }),
}));

type batchDataStoreType = {
  data: BatchData[];
  setData: (data: BatchData[]) => void;
};

export const useBatchDataStore = create<batchDataStoreType>((set) => ({
  data: [],
  setData: (data: BatchData[]) => set({ data: data }),
}));
