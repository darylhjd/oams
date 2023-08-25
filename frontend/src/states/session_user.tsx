import { User } from "@/api/models";
import { create } from "zustand";

export type sessionUserStoreType = {
  user: User | null;
  setUser: (user: User | null) => void;
}

export const sessionUserStore = create<sessionUserStoreType>((set) => ({
  user: null,
  setUser: (user: User | null) => {set({user: user})},
}))
