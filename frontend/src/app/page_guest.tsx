import styles from "@/styles/PageGuest.module.css";

import { Title } from "@mantine/core";

export default function GuestPage() {
  return (
    <Title ta="center" className={styles.title}>
      Welcome to the Online Attendance Management System
    </Title>
  );
}
