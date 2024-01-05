"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import NotFoundPage from "@/app/not-found";
import React from "react";

export function CheckHasUserRole({
  role,
  children,
}: {
  role: UserRole;
  children: React.ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.user.role != role) {
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

export function CanTakeAttendance({
  children,
  failNode,
}: {
  children: React.ReactNode;
  failNode: React.ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.management_details.attendance) {
    return <>{children}</>;
  }

  return <>{failNode}</>;
}

export function CanManageClassRulesAndReports({
  children,
  failNode,
}: {
  children: React.ReactNode;
  failNode: React.ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.management_details.rules_and_reports) {
    return <>{children}</>;
  }

  return <>{failNode}</>;
}
