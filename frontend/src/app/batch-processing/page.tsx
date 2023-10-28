"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import NotFoundPage from "../not-found";
import { Center } from "@mantine/core";

export default function BatchProcessingPage() {
  const session = useSessionUserStore();

  if (session.data?.session_user.role != UserRole.SystemAdmin) {
    return <NotFoundPage />;
  }

  return (
    <Center>
      <p>This is the batch processing page.</p>
    </Center>
  );
}
