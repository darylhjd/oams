"use client";

import { APIClient } from "@/api/client";
import { sessionStore } from "@/states/session";
import { Center, Container, Flex, Transition, createStyles } from "@mantine/core";
import { useEffect } from "react";

const useStyles = createStyles((theme) => ({
  container: {
    margin: "45vh 0",
  }
}))

export default function SessionInitialiser({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = sessionStore();
  const { classes } = useStyles();

  useEffect(() => {
    if (session.loaded) {
      return;
    }

    APIClient.getUserMe().then((data) => {
      session.setUser(data);
    });
  }, [session]);

  if (!session.loaded) {
    return (
      <Center>
        <Container className={classes.container}>
          <Flex align="center">
            <p>Beaming you to the site...</p>
          </Flex>
        </Container>        
      </Center>
    );
  }

  return (
    <Transition mounted transition="fade" duration={400} timingFunction="ease">
      {(styles) => <div style={styles}>{children}</div>}
    </Transition>
  );
}
