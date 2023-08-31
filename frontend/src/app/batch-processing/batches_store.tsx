import { BatchPostResponse } from "@/api/models";
import { create } from "zustand";

type batchesStoreType = {
  data: BatchPostResponse | null;
  setData: (batchResponse: BatchPostResponse | null) => void;
  invalidate: () => void;
};

export const batchesStore = create<batchesStoreType>((set) => ({
  data: null,
  setData: (batchResponse: BatchPostResponse | null) =>
    set({ data: batchResponse }),
  invalidate: () => set({ data: null }),
}));
