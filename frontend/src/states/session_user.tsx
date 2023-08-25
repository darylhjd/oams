import { User } from "@/api/models";
import { create } from "zustand";

export type sessionUserStoreType = {
  user: User | null;
  loaded: boolean;
  setUser: (user: User | null) => void;
}

export const sessionUserStore = create<sessionUserStoreType>((set) => ({
  user: null,
  loaded: false,
  setUser: (user: User | null) => {set({user: user, loaded: true})},
}))
