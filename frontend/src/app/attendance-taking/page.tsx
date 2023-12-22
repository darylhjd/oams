"use client";

import styles from "@/styles/AttendanceTakingPage.module.css";

import { APIClient } from "@/api/client";
import { Container, Group, Space, Text, Title, Tooltip } from "@mantine/core";
import { IconHelp } from "@tabler/icons-react";
import { useState } from "react";
import { UpcomingClassGroupSession } from "@/api/attendance_taking";
import { RequestLoader } from "@/components/request_loader";
import { SessionCard } from "@/app/attendance-taking/session_card";

export default function AttendanceTakingPage() {
  const [upcomingSessions, setUpcomingSessions] = useState<
    UpcomingClassGroupSession[]
  >([]);
  const promiseFunc = async () => {
    const data = await APIClient.attendanceTakingsGet();
    return setUpcomingSessions(data.upcoming_class_group_sessions);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <Container className={styles.container} fluid>
        <PageTitle />
        <UpcomingSessionsGrid sessions={upcomingSessions} />
      </Container>
    </RequestLoader>
  );
}

function PageTitle() {
  return (
    <Group className={styles.title} gap="xs" justify="center">
      <Title order={3} ta="center">
        Upcoming Class Group Sessions
      </Title>
      <Tooltip
        label="Only sessions beginning in less than 15 minutes are shown."
        events={{
          hover: true,
          focus: false,
          touch: true,
        }}
      >
        <IconHelp size={15} />
      </Tooltip>
    </Group>
  );
}

function UpcomingSessionsGrid({
  sessions,
}: {
  sessions: UpcomingClassGroupSession[];
}) {
  if (sessions.length == 0) {
    return (
      <Text ta="center">Hooray, you have no upcoming sessions to manage!</Text>
    );
  }

  const sessionCards = sessions.map((session) => (
    <SessionCard key={session.id} session={session} />
  ));

  return (
    <>
      <Text ta="center">
        You have{" "}
        <Text span c="green">
          {sessions.length}
        </Text>{" "}
        upcoming sessions.
      </Text>
      <Space h="md" />
      <Group justify="center">{sessionCards}</Group>
    </>
  );
}
