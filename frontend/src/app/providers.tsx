import { SessionInitialiser } from "@/stores/session";
import { MantineProvider } from "@mantine/core";
import { Notifications } from "@mantine/notifications";

export default function Providers({ children }: { children: React.ReactNode }) {
  return (
    <MantineProvider defaultColorScheme="dark">
      <Notifications />
      <SessionInitialiser>{children}</SessionInitialiser>
    </MantineProvider>
  );
}
