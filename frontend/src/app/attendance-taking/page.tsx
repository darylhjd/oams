"use client";

import styles from "@/styles/AttendanceTakingPage.module.css";

import { APIClient } from "@/api/client";
import {
  Container,
  Group,
  Paper,
  Space,
  Text,
  Title,
  Tooltip,
} from "@mantine/core";
import { IconHelp } from "@tabler/icons-react";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { Routes } from "@/routing/routes";
import { UpcomingClassGroupSession } from "@/api/attendance_taking";
import { UserRole } from "@/api/user";
import { RequestLoader } from "@/components/request_loader";

export default function AttendanceTakingPage() {
  const [upcomingSessions, setUpcomingSessions] = useState<
    UpcomingClassGroupSession[]
  >([]);
  const promiseFunc = async () => {
    const data = await APIClient.attendanceTakingGet();
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
      <>
        <Title visibleFrom="md" order={3}>
          Upcoming Class Group Sessions
        </Title>
        <Title hiddenFrom="md" order={4} ta="center">
          Upcoming Class Group Sessions
        </Title>
      </>

      <Tooltip
        label="Only sessions beginning in less than 15 minutes are shown."
        events={{
          hover: true,
          focus: false,
          touch: true,
        }}
      >
        <IconHelp size={12} />
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

function SessionCard({ session }: { session: UpcomingClassGroupSession }) {
  const router = useRouter();

  const startDatetime = new Date(session.start_time);
  const endDatetime = new Date(session.end_time);

  const date = startDatetime.toLocaleString(undefined, {
    day: "numeric",
    month: "numeric",
    year: "numeric",
  });
  const startTime = startDatetime.toLocaleString(undefined, {
    hour: "2-digit",
    minute: "2-digit",
  });
  const endTime = endDatetime.toLocaleString(undefined, {
    hour: "2-digit",
    minute: "2-digit",
  });

  const isOngoing = new Date() >= startDatetime;

  return (
    <Paper
      withBorder
      p="xs"
      className={styles.sessionCard}
      component="button"
      onClick={() => router.push(Routes.attendanceTakingSession + session.id)}
    >
      <Text ta="left">
        {session.code}{" "}
        <Text span size="sm" c="dimmed">
          {session.year}/{session.semester}
        </Text>
      </Text>
      <Text ta="left" size="xs">
        Class Type: {session.class_type}
      </Text>
      <Text ta="left" size="xs">
        Venue: {session.venue}
      </Text>
      <Space h="xs" />
      <Text ta="left" size="xs">
        {date}
        <br />
        {startTime} - {endTime}
      </Text>
      <Text ta="left" size="xs" c={isOngoing ? "green" : "orange"}>
        {isOngoing ? "ONGOING" : "STARTING"}
      </Text>
      <Space h="xs" />
      <Text ta="left" size="xs" c="dimmed">
        {session.managing_role ? session.managing_role : UserRole.SystemAdmin}
      </Text>
    </Paper>
  );
}
