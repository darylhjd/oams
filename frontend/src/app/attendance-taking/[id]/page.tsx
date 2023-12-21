"use client";

import styles from "@/styles/SessionAttendanceTaking.module.css";

import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { Center, Container, Text, Title } from "@mantine/core";
import {
  AttendanceTakingGetResponse,
  UpcomingClassGroupSession,
} from "@/api/attendance_taking";
import { SessionEnrollment } from "@/api/session_enrollment";

export default function SessionAttendanceTakingPage({
  params,
}: {
  params: Params;
}) {
  const [attendance, setAttendance] = useState<AttendanceTakingGetResponse>(
    {} as AttendanceTakingGetResponse,
  );
  const promiseFunc = async () => {
    const data = await APIClient.attendanceTakingGet(params.id);
    return setAttendance(data);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <PageTitle />
        <SessionInfo session={attendance.upcoming_class_group_session} />
        <AttendanceTaker enrollments={attendance.enrollment_data} />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Center className={styles.title}>
      <Title visibleFrom="md" order={3}>
        Attendance Taking
      </Title>
      <Title hiddenFrom="md" order={4}>
        Attendance Taking
      </Title>
    </Center>
  );
}

function SessionInfo({ session }: { session: UpcomingClassGroupSession }) {
  return null;
}

function AttendanceTaker({
  enrollments,
}: {
  enrollments: SessionEnrollment[];
}) {
  return (
    <Text ta="center">
      There are {enrollments.length} enrollment data rows.
    </Text>
  );
}
