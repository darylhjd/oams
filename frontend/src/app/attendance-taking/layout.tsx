import { CanTakeAttendance } from "@/components/session_checker";
import React from "react";
import NotFoundPage from "@/app/not-found";

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
  return (
    <CanTakeAttendance failNode={<NotFoundPage />}>
      {children}
    </CanTakeAttendance>
  );
}
