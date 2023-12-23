import { Session } from "@/api/session";
import { create } from "zustand";

type sessionUserStoreType = {
  data: Session | null;
  setSession: (data: Session | null) => void;
};

export const useSessionUserStore = create<sessionUserStoreType>((set) => ({
  data: null,
  setSession: (data: Session | null) => set({ data: data }),
}));
