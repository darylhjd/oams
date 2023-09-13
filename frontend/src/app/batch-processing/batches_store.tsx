import { BatchData, BatchPostResponse } from "@/api/models";
import { create } from "zustand";

type batchesStoreType = {
  data: BatchData[] | null;
  setData: (batchResponse: BatchPostResponse | null) => void;
  invalidate: () => void;
};

export const batchesStore = create<batchesStoreType>((set) => ({
  data: null,
  setData: (batchResponse: BatchPostResponse | null) => {
    if (batchResponse == null) {
      return;
    }

    var classGroupRunningIndex = 0;
    batchResponse.batches.forEach((batch, classIndex) => {
      batch.class_groups.forEach((classGroup) => {
        classGroup.class_id = classIndex;
        classGroup.sessions.forEach((session) => {
          session.class_group_id = classGroupRunningIndex;
        });
        classGroupRunningIndex++;
      });
    });

    set({ data: batchResponse.batches });
  },
  invalidate: () => set({ data: null }),
}));
