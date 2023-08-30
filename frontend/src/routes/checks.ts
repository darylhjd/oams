"use client"

import { sessionStore } from "@/states/session";
import { getURL } from "next/dist/shared/lib/utils";
import { Routes } from "./routes";
import { UserRole } from "@/api/models";

// Checks if there is already a user session. If yes, then redirect to the given URL.
export function redirectIfAlreadyLoggedIn(url: string): boolean {
  const session = sessionStore();

  if (session.data != null) {
    location.href = url;
    return true;
  }

  return false;
}

// Checks if there is a user session. If there is not, then redirect the user to login. This function
// automatically redirects the user back to the original URL.
export function redirectIfNotLoggedIn(): boolean {
  const session = sessionStore();
  const path = getURL();

  if (session.data == null) {
    const redirectUrl = `${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}${path}`;
    const queryParams = new URLSearchParams({
      redirect_url: redirectUrl,
    });

    location.replace(`${Routes.login}?${queryParams.toString()}`);
    return true;
  }

  return false;
}

// Checks if a user has required user role. If not, it redirects to home screen. Also checks if there is a user session.
// If not, it will redirect the user to log in first.
export function redirectIfNotUserRole(role: UserRole): boolean {
  const session = sessionStore();

  if (redirectIfNotLoggedIn()) {
    return true;
  }

  if (session.data!.session_user.role != role) {
    location.replace(Routes.home);
    return true;
  }

  return false;
}
