import { CanTakeAttendance } from "@/components/session_checker";
import NotFoundPage from "@/app/not-found";
import { ReactNode } from "react";

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
  children: ReactNode;
}) {
  return (
    <CanTakeAttendance failNode={<NotFoundPage />}>
      {children}
    </CanTakeAttendance>
  );
}
