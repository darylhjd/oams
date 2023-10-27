"use client";

import { APIClient } from "@/api/client";
import { sessionStore } from "@/states/session";
import { Center, Container, Flex, createStyles } from "@mantine/core";
import { useEffect, useState } from "react";

const useStyles = createStyles((theme) => ({
  container: {
    margin: "45vh 0",
  },
}));

export default function SessionInitialiser({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = sessionStore();
  const { classes } = useStyles();

  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    if (loaded) {
      return;
    }

    APIClient.getUserMe().then((data) => {
      session.setUser(data);
      setLoaded(true);
    });
  }, [session]);

  if (!loaded) {
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

  return <>{children}</>;
}
