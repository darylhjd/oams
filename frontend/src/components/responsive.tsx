"use client";

import { useMediaQuery } from "@mantine/hooks";

// A simple boolean flag to check if the current screen size is suitable for a
// mobile version or not.
export function isMobile() {
  return !useMediaQuery("(min-width: 62em)");
}
