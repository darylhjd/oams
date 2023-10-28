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
  setSession: (data: UserMeResponse) => void;
  invalidate: () => void;
};

export const useSessionUserStore = create<sessionUserStoreType>((set) => ({
  data: null,
  setSession: (data: UserMeResponse) => set({ data: data }),
  invalidate: () => set({ data: null }),
}));

export function SessionInitialiser({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = useSessionUserStore();
  const [loaded, setLoaded] = useState(false);

  useEffect(() => {
    if (loaded) {
      return;
    }

    APIClient.userMe()
      .then((data) => {
        session.setSession(data);
      })
      .catch((_error) => null)
      .finally(() => setLoaded(true));
  });

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
