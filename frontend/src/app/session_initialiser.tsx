'use client'

import { APIClient } from "@/api/client";
import { User } from "@/api/models";
import { sessionUserStore } from "@/states/session_user";
import { Center } from "@mantine/core";
import { useEffect } from "react";

export default function SessionInitialiser({children}: {children: React.ReactNode}) {
  const sessionStore = sessionUserStore()

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
