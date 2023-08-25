'use client'

import { APIClient } from "@/api/client";
import { userSessionStore } from "@/states/session_user";
import { Center } from "@mantine/core";
import { useEffect } from "react";

export default function SessionInitialiser({children}: {children: React.ReactNode}) {
  const sessionStore = userSessionStore()

  useEffect(() => {
    if (sessionStore.loaded) {
      return
    }

    APIClient.getUserMe().then((data) => {
      sessionStore.setUser(data)
    })
  })

  if (!sessionStore.loaded) {
    return (
      <Center>
        Loading...
      </Center>
    )
  }

  return (
    <>
      {children}
    </>
  )
}
