'use client'

import { sessionStore } from "@/states/session"
import LoggedHomePage from "./logged_page"
import GuestHomePage from "./guest_page"

export default function HomePage() {
  const session = sessionStore()

  switch (session.userMe) {
    case null:
      return <GuestHomePage />
    default:
      return <LoggedHomePage />
  }
}
