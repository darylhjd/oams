import {
  CanManageClassRules,
  CanTakeAttendance,
} from "@/components/session_checker";
import React from "react";
import NotFoundPage from "@/app/not-found";

export const metadata = {
  title: "Attendance Rules",
  description: "OAMS Attendance Rule Management Service",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function AttendanceRulesLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <CanManageClassRules failNode={<NotFoundPage />}>
      {children}
    </CanManageClassRules>
  );
}
