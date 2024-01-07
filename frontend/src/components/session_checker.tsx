"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import NotFoundPage from "@/app/not-found";
import { ReactNode } from "react";

export function CheckHasUserRole({
  role,
  children,
}: {
  role: UserRole;
  children: ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.user.role != role) {
    return <NotFoundPage />;
  }

  return <>{children}</>;
}

export function IsLoggedIn({ children }: { children: ReactNode }) {
  const session = useSessionUserStore();

  if (!session.data) {
    return <NotFoundPage />;
  }

  return <>{children}</>;
}

export function CanTakeAttendance({
  children,
  failNode,
}: {
  children: ReactNode;
  failNode: ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.management_details.attendance) {
    return <>{children}</>;
  }

  return <>{failNode}</>;
}

export function CanAdministerClass({
  children,
  failNode,
}: {
  children: ReactNode;
  failNode: ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.management_details.administrative) {
    return <>{children}</>;
  }

  return <>{failNode}</>;
}
