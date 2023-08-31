import { UserMeResponse } from "@/api/models";
import { create } from "zustand";

type sessionStoreType = {
  data: UserMeResponse | null;
  loaded: boolean;
  setUser: (user: UserMeResponse | null) => void;
  invalidate: () => void;
};

export const sessionStore = create<sessionStoreType>((set) => ({
  data: null,
  loaded: false,
  setUser: (user: UserMeResponse | null) => set({ data: user, loaded: true }),
  invalidate: () => set({ data: null, loaded: false }),
}));
