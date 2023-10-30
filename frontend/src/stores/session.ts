import { UserMeResponse } from "@/api/user";
import { create } from "zustand";

type sessionUserStoreType = {
  data: UserMeResponse | null;
  setSession: (data: UserMeResponse) => void;
};

export const useSessionUserStore = create<sessionUserStoreType>((set) => ({
  data: null,
  setSession: (data: UserMeResponse) => set({ data: data }),
}));
