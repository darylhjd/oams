// Import styles of packages that you've installed.
// All packages except `@mantine/hooks` require styles imports
import "@mantine/core/styles.css";
import "@mantine/dropzone/styles.css";
import "@mantine/notifications/styles.css";

import { ColorSchemeScript } from "@mantine/core";
import Providers from "./providers";
import Header from "@/components/header";
import CenteredPage from "@/components/centered_page";

export const metadata = {
  title: "OAMS",
  description: "Online Attendance Management System",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <head>
        <ColorSchemeScript defaultColorScheme="dark" />
      </head>
      <body>
        <Providers>
          <Header />
          <CenteredPage>{children}</CenteredPage>
        </Providers>
      </body>
    </html>
  );
}
