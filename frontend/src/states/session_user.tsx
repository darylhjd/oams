import { User } from "@/api/models";
import { create } from "zustand";

type userSessionStoreType = {
  user: User | null;
  loaded: boolean;
  setUser: (user: User | null) => void;
}

export const userSessionStore = create<userSessionStoreType>((set) => ({
  user: null,
  loaded: false,
  setUser: (user: User | null) => {set({user: user, loaded: true})},
}))
