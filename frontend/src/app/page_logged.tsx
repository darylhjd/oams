"use client";

import styles from "@/styles/PageLogged.module.css";

import { useSessionUserStore } from "@/stores/session";
import { Title } from "@mantine/core";

export default function LoggedPage() {
  const session = useSessionUserStore();

  return (
    <Title order={2} ta="center" className={styles.title}>
      Welcome, {session.data!.session_user.name}
    </Title>
  );
}
