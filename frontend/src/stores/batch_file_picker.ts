import { FileWithPath } from "@mantine/dropzone";
import { create } from "zustand";

export type batchFileStoreType = {
  files: FileWithPath[];
  setFiles: (file: FileWithPath[]) => void;
  clearFiles: () => void;
};

export const useBatchFiles = create<batchFileStoreType>((set) => ({
  files: [],
  setFiles: (files) => set({ files: files }),
  clearFiles: () => set({ files: [] }),
}));
