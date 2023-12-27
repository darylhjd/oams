"use client";

import styles from "@/styles/AttendanceRules.module.css";

import { Container, Title } from "@mantine/core";
import { useState } from "react";
import { CoordinatingClass } from "@/api/attendance_rule";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";

export default function AttendanceRulesPage() {
  const [coordinatingClasses, setCoordinatingClasses] = useState<
    CoordinatingClass[]
  >([]);
  const promiseFunc = async () => {
    const data = await APIClient.attendanceRulesGet();
    return setCoordinatingClasses(data.coordinating_classes);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <PageTitle />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Title order={3} ta="center">
      Attendance Rules
    </Title>
  );
}
