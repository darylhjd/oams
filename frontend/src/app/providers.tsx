import { SessionInitialiser } from "@/stores/session";
import { MantineProvider } from "@mantine/core";

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <MantineProvider defaultColorScheme="dark">
      <SessionInitialiser>{children}</SessionInitialiser>
    </MantineProvider>
  );
}
