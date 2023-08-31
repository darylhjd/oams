import { BatchPostResponse } from "@/api/models";
import { create } from "zustand";

type batchStoreType = {
  data: BatchPostResponse | null;
  setData: (batchResponse: BatchPostResponse | null) => void;
  invalidate: () => void;
};

export const batchStore = create<batchStoreType>((set) => ({
  data: null,
  setData: (batchResponse: BatchPostResponse | null) =>
    set({ data: batchResponse }),
  invalidate: () => set({ data: null }),
}));
