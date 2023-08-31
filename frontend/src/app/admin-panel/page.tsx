"use client";

import { UserRole } from "@/api/models";
import { buildRedirectUrlQueryParamsString } from "@/routes/checks";
import { Routes } from "@/routes/routes";
import { sessionStore } from "@/states/session";
import { Center } from "@mantine/core";
import { getURL } from "next/dist/shared/lib/utils";
import { useEffect } from "react";

export default function AdminPanelPage() {
  const session = sessionStore();

  useEffect(() => {
    if (session.data == null) {
      location.replace(
        `${Routes.login}?${buildRedirectUrlQueryParamsString(getURL())}`,
      );
    } else if (session.data!.session_user.role != UserRole.SystemAdmin) {
      location.replace(Routes.home);
    }
  }, [session]);

  if (
    session.data == null ||
    session.data!.session_user.role != UserRole.SystemAdmin
  ) {
    return null;
  }

  return (
    <Center>
      <p>This is the admin panel page.</p>
    </Center>
  );
}
