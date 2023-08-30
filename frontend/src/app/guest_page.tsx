"use client";

import { isMobile } from "@/components/responsive";
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
      <Title order={isMobile() ? 2 : 1} ta="center">
        Welcome to the Online Attendance Management System
      </Title>
    </div>
  );
}
