"use client";

import { useSessionUserStore } from "@/stores/session";
import { Center } from "@mantine/core";

export default function Home() {
  const session = useSessionUserStore();

  return (
    <Center>
      <p>Hello there, {session.data?.session_user.name}</p>
    </Center>
  );
}
