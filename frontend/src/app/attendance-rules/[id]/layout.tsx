import { Metadata } from "next";
import React from "react";
import { CanManageClassRules } from "@/components/session_checker";
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
    title: `Attendance Rules: ${params.id}`,
    description: "OAMS Attendance Rules",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function AttendanceRuleLayout({
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
