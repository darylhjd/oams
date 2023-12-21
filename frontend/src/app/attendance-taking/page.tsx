"use client";

import styles from "@/styles/AttendanceTakingPage.module.css";

import { ClassGroupSession } from "@/api/class_group_session";
import { APIClient } from "@/api/client";
import { IS_MOBILE_MEDIA_QUERY } from "@/components/media_query";
import {
  Center,
  Container,
  Group,
  Loader,
  Paper,
  SimpleGrid,
  Text,
  Title,
  Tooltip,
} from "@mantine/core";
import { useMediaQuery } from "@mantine/hooks";
import { IconHelp } from "@tabler/icons-react";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Routes } from "@/routing/routes";

export default function AttendanceTakingPage() {
  const [upcomingSessions, setUpcomingSessions] = useState<ClassGroupSession[]>(
    [],
  );
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    APIClient.attendanceTakingGet().then((data) => {
      setUpcomingSessions(data.upcoming_class_group_sessions);
      setLoaded(true);
    });
  }, []);

  if (!loaded) {
    return (
      <Center>
        <Loader />
      </Center>
    );
  }

  return (
    <Container className={styles.container} fluid>
      <PageTitle />
      <UpcomingSessionsGrid sessions={upcomingSessions} />
    </Container>
  );
}

function PageTitle() {
  const isMobile = useMediaQuery(IS_MOBILE_MEDIA_QUERY);

  return (
    <Group
      className={styles.title}
      gap="xs"
      justify={isMobile ? "flex-start" : "center"}
    >
      <>
        <Title visibleFrom="md" order={3}>
          Upcoming Class Group Sessions
        </Title>
        <Title hiddenFrom="md" order={4} ta="center">
          Upcoming Class Group Sessions
        </Title>
      </>

      <Tooltip
        label="Only sessions starting in 15 minutes will be shown!"
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

function UpcomingSessionsGrid({ sessions }: { sessions: ClassGroupSession[] }) {
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
      <Text ta="center">You have {sessions.length} upcoming sessions.</Text>
      <>
        <SimpleGrid visibleFrom="sm" cols={3}>
          {sessionCards}
        </SimpleGrid>
        <SimpleGrid hiddenFrom="sm" cols={1}>
          {sessionCards}
        </SimpleGrid>
      </>
    </>
  );
}

function SessionCard({ session }: { session: ClassGroupSession }) {
  const router = useRouter();

  return (
    <Paper
      withBorder
      className={styles.sessionCard}
      component="button"
      onClick={() => router.push(Routes.attendanceTakingSession + session.id)}
    >
      <Text ta="center">{session.venue}</Text>
    </Paper>
  );
}
