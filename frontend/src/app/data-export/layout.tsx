import { UserRole } from "@/api/user";
import { CheckHasUserRole } from "@/components/session_checker";
import React from "react";

export const metadata = {
  title: "Data Export",
  description: "OAMS Data Export Service",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function DataExportLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <CheckHasUserRole role={UserRole.SystemAdmin}>{children}</CheckHasUserRole>
  );
}
