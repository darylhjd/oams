import { create } from "zustand";
import { FileStoreType } from "@/stores/file_store";
import { UpsertClassGroupManagerParams } from "@/api/class_group_manager";

export const useManagerFilesStore = create<FileStoreType>((set) => ({
  files: [],
  setFiles: (files) => set({ files: files }),
  clearFiles: () => set({ files: [] }),
}));

type managerDataStoreType = {
  data: UpsertClassGroupManagerParams[];
  setData: (data: UpsertClassGroupManagerParams[]) => void;
};

export const useManagerDataStore = create<managerDataStoreType>((set) => ({
  data: [],
  setData: (data: UpsertClassGroupManagerParams[]) => set({ data: data }),
}));
