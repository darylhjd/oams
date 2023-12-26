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

export function HasManagedClassGroups({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = useSessionUserStore();

  if (session.data?.management_details.has_managed_class_groups) {
    return <>{children}</>;
  }

  return <NotFoundPage />;
}
