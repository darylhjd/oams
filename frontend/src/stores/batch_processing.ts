import { BatchData } from "@/api/batch";
import { FileWithPath } from "@mantine/dropzone";
import { create } from "zustand";

type batchFileStoreType = {
  files: FileWithPath[];
  setFiles: (file: FileWithPath[]) => void;
  clearFiles: () => void;
};

export const useBatchFilesStore = create<batchFileStoreType>((set) => ({
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
