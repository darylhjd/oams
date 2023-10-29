"use client";

import styles from "@/styles/Providers.module.css";

import { APIClient } from "@/api/client";
import { Center, MantineProvider, Skeleton } from "@mantine/core";
import { Notifications } from "@mantine/notifications";
import { motion } from "framer-motion";
import { useSessionUserStore } from "@/stores/session";
import { useEffect, useState } from "react";

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <MantineProvider defaultColorScheme="dark">
      <Notifications />
      <SessionInitialiser>{children}</SessionInitialiser>
    </MantineProvider>
  );
}

function SessionInitialiser({ children }: { children: React.ReactNode }) {
  const session = useSessionUserStore();
  const [loaded, setLoaded] = useState(false);
  const [fetching, setFetching] = useState(false);

  useEffect(() => {
    if (fetching) {
      return;
    }

    setFetching(true);
    APIClient.userMe()
      .then((data) => {
        session.setSession(data);
      })
      .catch((_error) => null)
      .finally(() => setLoaded(true));
  }, [fetching, session]);

  return loaded ? (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      transition={{ duration: 1 }}
    >
      {children}
    </motion.div>
  ) : (
    <Center className={styles.center}>
      <Skeleton className={styles.skeleton} visible={true} />
    </Center>
  );
}