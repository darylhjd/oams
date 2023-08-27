import { UserMeResponse } from "@/api/models";
import { create } from "zustand";

type sessionStoreType = {
  userMe: UserMeResponse | null;
  loaded: boolean;
  setUser: (user: UserMeResponse | null) => void;
  invalidate: () => void;
};

export const sessionStore = create<sessionStoreType>((set) => ({
  userMe: null,
  loaded: false,
  setUser: (user: UserMeResponse | null) => set({ userMe: user, loaded: true }),
  invalidate: () => set({ userMe: null, loaded: false }),
}));
