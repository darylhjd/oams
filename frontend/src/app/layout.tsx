// Import styles of packages that you've installed.
// All packages except `@mantine/hooks` require styles imports
import "@mantine/core/styles.css";
import "@mantine/code-highlight/styles.css";
import "@mantine/dropzone/styles.css";
import "@mantine/notifications/styles.css";
import "mantine-react-table/styles.css";

import { ColorSchemeScript } from "@mantine/core";
import Header from "@/components/header";
import CenteredPage from "@/components/centered_page";
import Providers from "@/components/providers";
import { ReactNode } from "react";

export const metadata = {
  title: "OAMS",
  description: "Online Attendance Management System",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function RootLayout({ children }: { children: ReactNode }) {
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
