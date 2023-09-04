import { sessionStore } from "@/states/session";
import { useRouter } from "next/navigation";
import { Routes } from "./routes";
import { getURL } from "next/dist/shared/lib/utils";
import { UserRole } from "@/api/models";

// Redirect the user back to the home screen if user is already logged in.
// This is useful, for example, in the event the user purposefully tries to enter the
// login page when they are already logged in.
export function redirectIfLoggedIn(): boolean {
  const router = useRouter();
  const session = sessionStore();

  if (session.data != null) {
    router.replace(Routes.home);
    return true;
  }

  return false;
}

// Redirect the user to the login page if they are not logged in. This will also help
// automatically set the redirect for the login page so that they are transported back
// to the original page after login.
export function redirectIfNotLoggedIn(): boolean {
  const router = useRouter();
  const session = sessionStore();

  if (session.data == null) {
    router.replace(loginWithRedirectUrl());
    return true;
  }

  return false;
}

// Redirects the user back to the home screen if they do not have the appropriate user role.
export function redirectIfNotUserRole(userRole: UserRole): boolean {
  const router = useRouter();
  const session = sessionStore();

  if (redirectIfNotLoggedIn()) {
    return true;
  }

  if (session.data!.session_user.role != userRole) {
    router.replace(Routes.home);
    return true;
  }

  return false;
}

function loginWithRedirectUrl(): string {
  const path = getURL();
  const redirectUrl = `${process.env.WEB_SERVER_HOST}:${process.env.WEB_SERVER_PORT}${path}`;
  const queryParams = new URLSearchParams({
    redirect_url: redirectUrl,
  });
  return `${Routes.login}?${queryParams.toString()}`;
}
