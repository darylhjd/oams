"use client";

import { useSessionUserStore } from "@/stores/session";
import { Text } from "@mantine/core";

export default function ProfilePage() {
  const session = useSessionUserStore();

  return <Text ta="center">{session.data?.session_user.name}</Text>;
}
