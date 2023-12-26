import styles from "@/styles/AttendanceRules.module.css";

import { Container, Title } from "@mantine/core";

export default function AttendanceRulesPage() {
  return (
    <Container className={styles.container} fluid>
      <PageTitle />
    </Container>
  );
}

function PageTitle() {
  return (
    <Title order={3} ta="center">
      Attendance Rules
    </Title>
  );
}
