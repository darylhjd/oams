import styles from "@/styles/ClassAdministrationPage.module.css";

import { ReactNode } from "react";

export function Panel({ children }: { children: ReactNode }) {
  return <div className={styles.tabPanel}>{children}</div>;
}
