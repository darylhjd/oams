import { Title } from "@mantine/core";
import styles from "@/styles/PageGuest.module.css";

export default function GuestPage() {
  return (
    <Title ta="center" className={styles.title}>
      Welcome to the Online Attendance Management System
    </Title>
  );
}
