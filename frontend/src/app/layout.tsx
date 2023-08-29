import type { Metadata } from "next";
import React from "react";
import Providers from "@/app/providers";
import Header from "@/components/header";
import CenteredScreen from "@/components/centered_page";
import SessionInitialiser from "./session_initialiser";

export const metadata: Metadata = {
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
      <body>
        <Providers>
          <Header />
          <CenteredScreen>{children}</CenteredScreen>
        </Providers>
      </body>
    </html>
  );
}
