'use client'

import { APIClient } from "@/api/client"
import { sessionUserStore } from "@/states/session_user"
import { Center } from "@mantine/core"
import { useEffect } from "react"

export default function Page() {
  const sessionUser = sessionUserStore()

  useEffect(() => {
    if (sessionUser.loaded) {
      return
    }

    APIClient.getUserMe().then((data) => {
      sessionUser.setUser(data)
    })
  })

  if (!sessionUser.loaded) {
    return (
      <Center>
        <p>Loading user....</p>
      </Center>
    )
  }

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
