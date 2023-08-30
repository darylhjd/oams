"use client";

import { Title, createStyles } from "@mantine/core";
import { useMediaQuery } from "@mantine/hooks";

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

  const largeScreen = useMediaQuery("(min-width: 62em)");

  return (
    <div className={classes.hero}>
      <Title order={largeScreen ? 1 : 2} ta="center">
        Welcome to the Online Attendance Management System
      </Title>
    </div>
  );
}
