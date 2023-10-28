"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import { Center } from "@mantine/core";
import NotFoundPage from "../not-found";

export default function AdminPanelPage() {
  const session = useSessionUserStore();

  if (session.data?.session_user.role != UserRole.SystemAdmin) {
    return <NotFoundPage />;
  }

  return (
    <Center>
      <p>This is the admin panel page.</p>
    </Center>
  );
}
