'use client'

import { sessionStore } from "@/states/session"
import { Center } from "@mantine/core"

export default function Page() {
  const session = sessionStore()

  return (
    <Center>
      <p>
        This is the home page.
        <br />
        Has user session: {`${!(session.user == null)}`}
      </p>
    </Center>
  )
}
