import styles from "@/styles/CenteredPage.module.css";
import React from "react";

export default function CenteredPage({
  children,
}: {
  children: React.ReactNode;
}) {
  return <div className={styles.centered}>{children}</div>;
}
