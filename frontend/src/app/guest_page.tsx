"use client";

import { MOBILE_MIN_WIDTH } from "@/components/responsive";
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

  return (
    <div className={classes.hero}>
      <Title order={useMediaQuery(MOBILE_MIN_WIDTH) ? 2 : 1} ta="center">
        Welcome to the Online Attendance Management System
      </Title>
    </div>
  );
}
