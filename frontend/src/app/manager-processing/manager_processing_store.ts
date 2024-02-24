import { create } from "zustand";
import { UpsertClassGroupManagerParams } from "@/api/class_group_manager";
import { FileSetter } from "@/components/file_processing";

export type ManagerFileStoreType = FileSetter;

export const useManagerFilesStore = create<ManagerFileStoreType>((set) => ({
  files: [],
  setFiles: (files) => set({ files: files }),
  resetFiles: () => set({ files: [] }),
}));

export type ManagerDataStoreType = {
  data: UpsertClassGroupManagerParams[];
  setData: (data: UpsertClassGroupManagerParams[]) => void;
};

export const useManagerDataStore = create<ManagerDataStoreType>((set) => ({
  data: [],
  setData: (data: UpsertClassGroupManagerParams[]) => set({ data: data }),
}));
