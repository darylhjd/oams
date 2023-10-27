"use client";

import { useSessionUserStore } from "@/stores/session";
import LoggedPage from "./page_logged";
import GuestPage from "./page_guest";

export default function Home() {
  const session = useSessionUserStore();

  return session.data ? LoggedPage() : GuestPage();
}
