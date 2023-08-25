'use client'

import { userSessionStore } from "@/states/session_user"
import { Center } from "@mantine/core"

export default function Page() {
  const sessionUser = userSessionStore()

  return (
    <Center>
      <p>
        This is the home page.
        <br />
        Has user session: {`${!(sessionUser.user == null)}`}
      </p>
    </Center>
  )
}
