'use client'

import { RedirectIfNotLoggedIn } from "@/routes/checks";
import { Center } from "@mantine/core";

export default function LoginPage() {
  RedirectIfNotLoggedIn()

  return (
    <Center>
      <p>This is the profile screen.</p>
    </Center>
  )
}
