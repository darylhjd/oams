import { CanManageClassRulesAndReports } from "@/components/session_checker";
import React from "react";
import NotFoundPage from "@/app/not-found";

export const metadata = {
  title: "Class Reports",
  description: "OAMS Class Reports Service",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function ClassReportsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <CanManageClassRulesAndReports failNode={<NotFoundPage />}>
      {children}
    </CanManageClassRulesAndReports>
  );
}
