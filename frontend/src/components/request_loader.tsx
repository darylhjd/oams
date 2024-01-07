import styles from "@/styles/EntityLoader.module.css";

import NotFoundPage from "@/app/not-found";
import { Center, Loader } from "@mantine/core";
import { isAxiosError } from "axios";
import { ReactNode, useEffect, useState } from "react";
import { notifications } from "@mantine/notifications";
import { IconX } from "@tabler/icons-react";
import { motion } from "framer-motion";

export function RequestLoader({
  promiseFunc,
  children,
}: {
  promiseFunc: () => Promise<void>;
  children: ReactNode;
}) {
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState<any>(null);

  useEffect(() => {
    promiseFunc()
      .catch((error) => setError(error))
      .finally(() => setLoaded(true));
  }, []);

  if (!loaded) {
    return (
      <Center>
        <Loader className={styles.loader} />
      </Center>
    );
  } else if (error) {
    if (!isAxiosError(error) || error.response?.status != 404) {
      notifications.show({
        title: "API Error!",
        message:
          "There was a problem while executing the API request. Try again later.",
        icon: <IconX />,
        color: "red",
      });
      return null;
    } else {
      return <NotFoundPage />;
    }
  } else {
    return (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.5 }}
      >
        {children}
      </motion.div>
    );
  }
}
