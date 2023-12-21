import { IsClassGroupManager } from "@/components/session_checker";
import React from "react";

export const metadata = {
  title: "Attendance Taking",
  description: "OAMS Attendance Taking Service",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function AttendanceTakingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <IsClassGroupManager>{children}</IsClassGroupManager>;
}
