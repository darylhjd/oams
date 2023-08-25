'use client'

import { APIClient } from "@/api/client"
import { sessionUserStore } from "@/states/session_user"
import { Center } from "@mantine/core"
import { useEffect, useState } from "react"

export default function Page() {
  const [hasLoaded, setLoaded] = useState(false)
  const sessionUser = sessionUserStore()

  useEffect(() => {
    if (hasLoaded) {
      return
    }

    APIClient.getUserMe().then((data) => {
      setLoaded(true)
      sessionUser.setUser(data)
    })
  })

  if (!hasLoaded) {
    return (
      <div>
        Loading user....
      </div>
    )
  }

  return (
    <div>
      <Center>
        This is the home page.
        <br />
        Has user session: {`${!(sessionUser.user == null)}`}
      </Center>
    </div>
  )
}
