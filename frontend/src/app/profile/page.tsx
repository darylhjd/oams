'use client'

import { redirectIfNotLoggedIn } from "@/routes/checks";
import { Center } from "@mantine/core";

export default function LoginPage() {
  redirectIfNotLoggedIn()

  return (
    <Center>
      <p>This is the profile screen.</p>
    </Center>
  )
}
