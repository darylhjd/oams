import { sessionStore } from "@/states/session";
import { Routes } from "./routes";
import { getURL } from "next/dist/shared/lib/utils";
import { UserRole } from "@/api/models";

// Redirect the user back to the home screen if user is already logged in.
// This is useful, for example, in the event the user purposefully tries to enter the
// login page when they are already logged in.
export function redirectIfLoggedIn(): boolean {
  const session = sessionStore();

  if (session.data != null) {
    location.replace(Routes.index);
    return true;
  }

  return false;
}

// Redirect the user to the login page if they are not logged in. This will also help
// automatically set the redirect for the login page so that they are transported back
// to the original page after login. isProtected will show a 404 screen instead of a redirect to
// the login, which is useful for not exposing private webpages.
export function redirectIfNotLoggedIn(isProtected: boolean = false): boolean {
  const session = sessionStore();

  if (session.data == null) {
    location.replace(isProtected ? Routes.notFound : loginWithRedirectUrl());
    return true;
  }

  return false;
}

// Redirects the user back to the home screen if they do not have the appropriate user role.
// Note that user role that is not User will cause the redirect to show the 404 screen instead
// if the login page.
export function redirectIfNotUserRole(userRole: UserRole): boolean {
  const session = sessionStore();

  if (redirectIfNotLoggedIn(userRole != UserRole.User)) {
    return true;
  }

  if (session.data!.session_user.role != userRole) {
    location.replace(Routes.index);
    return true;
  }

  return false;
}

// Returns a URL to the login page, with the redirect_url query parameter set to the current
// page before the redirect.
function loginWithRedirectUrl(): string {
  const path = getURL();
  const redirectUrl = `${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}${path}`;
  const queryParams = new URLSearchParams({
    redirect_url: redirectUrl,
  });
  return `${Routes.login}?${queryParams.toString()}`;
}
