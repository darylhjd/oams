import { FileWithPath } from "@mantine/dropzone";
import { create } from "zustand";

export type batchFileStoreType = {
  files: FileWithPath[];
  setFiles: (file: FileWithPath[]) => void;
  clearFiles: () => void;
};

export const useBatchFilesStore = create<batchFileStoreType>((set) => ({
  files: [],
  setFiles: (files) => set({ files: files }),
  clearFiles: () => set({ files: [] }),
}));
