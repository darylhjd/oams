'use client'

import { sessionUserStore } from "@/states/session_user"
import { Center } from "@mantine/core"

export default function Page() {
  const sessionUser = sessionUserStore()

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
