"use client";

import { UserRole } from "@/api/user";
import NotFoundPage from "@/app/not-found";
import { useSessionUserStore } from "@/stores/session";

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
