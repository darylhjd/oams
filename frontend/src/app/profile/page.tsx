"use client";

import styles from "@/styles/ProfilePage.module.css";

import { useSessionUserStore } from "@/stores/session";
import { Paper, Space, Text } from "@mantine/core";

export default function ProfilePage() {
  const session = useSessionUserStore();
  const user = session.data!.session_user;

  return (
    <Paper className={styles.paper} radius="md" shadow="xs" withBorder p="xl">
      <Text ta="center" size="xl" fw={1000}>
        {user.id}
      </Text>
      <Space h="md" />
      <Text ta="center" size="sm">
        {user.role} • {user.email}
      </Text>
      <Space h="xs" />
      <Text c="dimmed" fs="italic" ta="center" size="sm">
        {user.name ? user.name : "No name registed"}
      </Text>
    </Paper>
  );
}
