"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import NotFoundPage from "@/app/not-found";

export function CheckHasUserRole({
  role,
  children,
}: {
  role: UserRole;
  children: React.ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.session_user.role != role) {
    return <NotFoundPage />;
  }

  return <>{children}</>;
}

export function IsLoggedIn({ children }: { children: React.ReactNode }) {
  const session = useSessionUserStore();

  if (!session.data) {
    return <NotFoundPage />;
  }

  return <>{children}</>;
}

export function IsClassGroupManager({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = useSessionUserStore();

  if (
    session.data?.has_managed_class_groups ||
    session.data?.session_user.role == UserRole.SystemAdmin
  ) {
    return <>{children}</>;
  }

  return <NotFoundPage />;
}
