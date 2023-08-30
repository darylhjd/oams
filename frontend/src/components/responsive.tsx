import { MediaQuery } from "@mantine/core";
import { useMediaQuery } from "@mantine/hooks";

// Shows children in mobile mode.
export function Mobile({ children }: { children: React.ReactNode }) {
  return (
    <MediaQuery largerThan="md" styles={{ display: "none" }}>
      {children}
    </MediaQuery>
  );
}

// Shows children in desktop mode.
export function Desktop({ children }: { children: React.ReactNode }) {
  return (
    <MediaQuery smallerThan="md" styles={{ display: "none" }}>
      {children}
    </MediaQuery>
  );
}

// A simple boolean flag to check if the current screen size is suitable for a
// desktop version or not.
export function isDesktop() {
  return useMediaQuery("(min-width: 62em)");
}
