"use client";

import { isDesktop } from "@/components/responsive";
import { Title, createStyles } from "@mantine/core";

const useStyles = createStyles((theme) => ({
  hero: {
    padding: "3em 0",
  },
}));

export default function GuestHomePage() {
  return <TopBanner />;
}

function TopBanner() {
  const { classes } = useStyles();

  return (
    <div className={classes.hero}>
      <Title order={isDesktop() ? 1 : 2} ta="center">
        Welcome to the Online Attendance Management System
      </Title>
    </div>
  );
}
