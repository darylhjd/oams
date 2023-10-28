"use client";

import { APIClient } from "@/api/client";
import { UserMeResponse } from "@/api/user";
import { Center, Skeleton } from "@mantine/core";
import { useEffect, useState } from "react";
import { create } from "zustand";
import { motion } from "framer-motion";
import styles from "@/styles/SessionInitialiser.module.css";

type sessionUserStoreType = {
  data: UserMeResponse | null;
  loaded: boolean;
  setSession: (data: UserMeResponse) => void;
  setLoaded: () => void;
};

export const useSessionUserStore = create<sessionUserStoreType>((set) => ({
  data: null,
  loaded: false,
  setSession: (data: UserMeResponse) => set({ data: data, loaded: true }),
  setLoaded: () => set({ loaded: true }),
}));

export function SessionInitialiser({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = useSessionUserStore();
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
      .finally(() => session.setLoaded());
  }, [fetching, session]);

  return session.loaded ? (
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
