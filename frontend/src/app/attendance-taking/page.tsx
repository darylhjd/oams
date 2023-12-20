import styles from "@/styles/AttendanceTakingPage.module.css";

import { Container, Text } from "@mantine/core";

export default function AttendanceTakingPage() {
  return (
    <Container className={styles.container} fluid>
      <Text ta="center">This is the attendance taking page.</Text>
    </Container>
  );
}
