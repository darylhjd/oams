import { User } from "@/api/models";
import { create } from "zustand";

type sessionStoreType = {
  user: User | null;
  loaded: boolean;
  setUser: (user: User | null) => void;
}

export const sessionStore = create<sessionStoreType>((set) => ({
  user: null,
  loaded: false,
  setUser: (user: User | null) => set({user: user, loaded: true}),
}))
