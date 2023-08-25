'use client'

import { APIClient } from "@/api/client";
import { sessionStore } from "@/states/session";
import { Center } from "@mantine/core";
import { useEffect } from "react";

export default function SessionInitialiser({children}: {children: React.ReactNode}) {
  const session = sessionStore()

  useEffect(() => {
    if (session.loaded) {
      return
    }

    APIClient.getUserMe().then((data) => {
      session.setUser(data)
    })
  }, [])

  if (!session.loaded) {
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
