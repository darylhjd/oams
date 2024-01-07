import styles from "@/styles/CenteredPage.module.css";
import { ReactNode } from "react";

export default function CenteredPage({ children }: { children: ReactNode }) {
  return <div className={styles.centered}>{children}</div>;
}
