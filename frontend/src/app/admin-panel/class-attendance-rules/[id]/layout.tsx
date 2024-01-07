import { Metadata } from "next";
import { ReactNode } from "react";

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
  children: ReactNode;
}) {
  return <>{children}</>;
}
