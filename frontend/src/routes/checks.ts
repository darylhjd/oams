import { sessionStore } from "@/states/session"
import { getURL } from "next/dist/shared/lib/utils"
import { redirect } from "next/navigation"
import { Routes } from "./routes"

// Checks if there is already a user session. If yes, then redirect to the given URL.
export function RedirectIfAlreadyLoggedIn(url: string) {
  const session = sessionStore()

  if (session.user != null) {
    redirect(url)
  }
}

// Checks if there is a user session. If there is not, then redirect the user to login. This function
// automatically redirects the user back to the original URL.
export function RedirectIfNotLoggedIn() {
  const session = sessionStore()
  const path = getURL()

  if (session.user == null) {
    const redirectUrl = `${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}${path}`
    const queryParams = new URLSearchParams({
      redirect_url: redirectUrl,
    })
    redirect(`${Routes.login}?${queryParams.toString()}`)
  }
}
