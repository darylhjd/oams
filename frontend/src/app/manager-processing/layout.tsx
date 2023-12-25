import { UserRole } from "@/api/user";
import { CheckHasUserRole } from "@/components/session_checker";
import React from "react";

export const metadata = {
  title: "Manager Processing",
  description: "OAMS Manager Processing Service",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function ManagerProcessingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <CheckHasUserRole role={UserRole.SystemAdmin}>{children}</CheckHasUserRole>
  );
}
