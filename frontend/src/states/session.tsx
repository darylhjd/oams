import { UserMeResponse } from "@/api/models";
import { create } from "zustand";

type sessionStoreType = {
  user: UserMeResponse | null;
  loaded: boolean;
  setUser: (user: UserMeResponse | null) => void;
  invalidate: () => void;
}

export const sessionStore = create<sessionStoreType>((set) => ({
  user: null,
  loaded: false,
  setUser: (user: UserMeResponse | null) => set({user: user, loaded: true}),
  invalidate: () => set({user: null, loaded: false})
}))
