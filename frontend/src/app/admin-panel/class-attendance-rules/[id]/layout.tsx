import { Metadata } from "next";
import React from "react";

export type Params = {
  id: number;
};

export async function generateMetadata({
  params,
}: {
  params: Params;
}): Promise<Metadata> {
  return {
    title: `Class Attendance Rule: ${params.id}`,
    description: "OAMS Class Attendance Rule",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function AdminPanelClassAttendanceRuleLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
