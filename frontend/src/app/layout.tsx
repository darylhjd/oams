import type { Metadata } from "next";
import React, { Suspense } from "react";
import Providers from "@/app/providers";
import Header from "@/components/header";
import CenteredScreen from "@/components/centered_page";

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
        <Suspense>
          <Providers>
            <Header />
            <CenteredScreen>{children}</CenteredScreen>
          </Providers>          
        </Suspense>
      </body>
    </html>
  );
}
