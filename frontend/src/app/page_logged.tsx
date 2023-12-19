"use client";

import styles from "@/styles/PageLogged.module.css";

import { useSessionUserStore } from "@/stores/session";
import { Title } from "@mantine/core";

export default function LoggedPage() {
  const session = useSessionUserStore();

  const name = session.data!.session_user.name;

  return (
    <Title order={2} ta="center" className={styles.title}>
      Welcome, {name == "" ? session.data!.session_user.email : name}
    </Title>
  );
}
