"use client";

import { UserRole } from "@/api/user";
import { useSessionUserStore } from "@/stores/session";
import { Text, Container } from "@mantine/core";
import NotFoundPage from "../not-found";
import styles from "@/styles/AdminPage.module.css";

export default function AdminPanelPage() {
  const session = useSessionUserStore();

  if (session.data?.session_user.role != UserRole.SystemAdmin) {
    return <NotFoundPage />;
  }

  return (
    <Container className={styles.container} fluid>
      <Text ta="center">This is the admin page!</Text>
    </Container>
  );
}
