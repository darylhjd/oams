"use client";

import { UserRole } from "@/api/models";
import { redirectIfNotUserRole } from "@/routes/checks";
import { Center } from "@mantine/core";

export default function BatchProcessingPage() {
  redirectIfNotUserRole(UserRole.SystemAdmin);

  return (
    <Center>
      <p>This is the batch processing page.</p>
    </Center>
  );
}
