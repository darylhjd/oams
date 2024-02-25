import { BatchData } from "@/api/batch";
import { create } from "zustand";
import { FileWithPath } from "@mantine/dropzone";
import { FileSetter } from "@/components/file_processing";

export const defaultSem1StartWeek = 33;
export const defaultSem2StartWeek = 2;

export type BatchFileStoreType = {
  startWeek: number;
  setStartWeek: (startWeek: number) => void;
  reset: () => void;
} & FileSetter;

// This is used to store the files and the start week that is sent to the POST
// endpoint.
export const useBatchFilesStore = create<BatchFileStoreType>((set) => ({
  files: [],
  setFiles: (files: FileWithPath[]) => set({ files }),
  resetFiles: () => set({ files: [] }),
  startWeek: defaultSem1StartWeek,
  setStartWeek: (startWeek: number) => set({ startWeek: startWeek }),
  reset: () => set({ files: [], startWeek: defaultSem1StartWeek }),
}));

export type BatchDataStoreType = {
  data: BatchData[];
  setData: (data: BatchData[]) => void;
};

// This is used to store the BatchData that is returned from POST and that is
// sent to PUT.
export const useBatchDataStore = create<BatchDataStoreType>((set) => ({
  data: [],
  setData: (data: BatchData[]) => set({ data: data }),
}));
