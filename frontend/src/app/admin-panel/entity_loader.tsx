import styles from "@/styles/EntityLoader.module.css";

import NotFoundPage from "@/app/not-found";
import { Center, Loader, Text } from "@mantine/core";
import { isAxiosError } from "axios";
import { useEffect, useState } from "react";

export function EntityLoader({
  promiseFunc,
  children,
}: {
  promiseFunc: () => Promise<void>;
  children: React.ReactNode;
}) {
  const [loaded, setLoaded] = useState(false);
  const [fetching, setFetching] = useState(false);
  const [error, setError] = useState<any | null>(null);

  useEffect(() => {
    if (fetching) {
      return;
    }

    setFetching(true);
    promiseFunc()
      .catch((error) => setError(error))
      .finally(() => setLoaded(true));
  }, [fetching, promiseFunc]);

  if (!loaded) {
    return (
      <Center>
        <Loader className={styles.loader} />
      </Center>
    );
  } else if (error) {
    if (!isAxiosError(error) || error.response?.status != 404) {
      return <Text ta="center">Error getting user!</Text>;
    } else {
      return <NotFoundPage />;
    }
  } else {
    return <>{children}</>;
  }
}
