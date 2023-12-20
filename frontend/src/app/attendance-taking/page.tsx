"use client";

import { ClassGroupSession } from "@/api/class_group_session";
import { APIClient } from "@/api/client";
import styles from "@/styles/AttendanceTakingPage.module.css";

import { Center, Container, Loader, Text } from "@mantine/core";
import { useEffect, useState } from "react";

export default function AttendanceTakingPage() {
  const [upcomingSessions, setUpcomingSessions] = useState<ClassGroupSession[]>(
    [],
  );
  const [loaded, setLoaded] = useState(false);
  const [fetching, setFetching] = useState(false);

  useEffect(() => {
    if (fetching) {
      return;
    }

    setFetching(true);
    APIClient.attendanceTakingGet().then((data) => {
      setUpcomingSessions(data.upcoming_class_group_sessions);
      setLoaded(true);
    });
  });

  if (!loaded) {
    return (
      <Center>
        <Loader />
      </Center>
    );
  }

  return (
    <Container className={styles.container} fluid>
      <Text ta="center">
        You have {upcomingSessions.length} upcoming attendance takings.
      </Text>
    </Container>
  );
}
