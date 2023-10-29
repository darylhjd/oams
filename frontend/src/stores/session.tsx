"use client";

import styles from "@/styles/SessionInitialiser.module.css";

import { APIClient } from "@/api/client";
import { UserMeResponse } from "@/api/user";
import { Center, Skeleton } from "@mantine/core";
import { useEffect, useState } from "react";
import { create } from "zustand";
import { motion } from "framer-motion";

type sessionUserStoreType = {
  data: UserMeResponse | null;
  setSession: (data: UserMeResponse) => void;
};

export const useSessionUserStore = create<sessionUserStoreType>((set) => ({
  data: null,
  setSession: (data: UserMeResponse) => set({ data: data }),
}));

export function SessionInitialiser({
  children,
}: {
  children: React.ReactNode;
}) {
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
