'use client'

import { UserRole } from "@/api/models";
import { redirectIfNotUserRole } from "@/routes/checks";
import { Center } from "@mantine/core";

export default function AdminPanelPage() {
  redirectIfNotUserRole(UserRole.SystemAdmin)

  return (
    <Center>
      <p>This is the admin panel page.</p> 
    </Center>
  )
}
