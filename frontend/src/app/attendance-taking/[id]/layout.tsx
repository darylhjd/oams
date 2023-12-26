import { Metadata } from "next";
import React from "react";
import { CanTakeAttendance } from "@/components/session_checker";
import NotFoundPage from "@/app/not-found";

export type Params = {
  id: number;
};

export async function generateMetadata({
  params,
}: {
  params: Params;
}): Promise<Metadata> {
  return {
    title: `Attendance Taking: ${params.id}`,
    description: "OAMS Attendance Taking",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function SessionAttendanceTakingLayout({
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
