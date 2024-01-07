import { CanAdministerClass } from "@/components/session_checker";
import NotFoundPage from "@/app/not-found";
import { ReactNode } from "react";

export const metadata = {
  title: "Class Administration",
  description: "OAMS Class Administration Service",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function ClassAdministrationsLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <CanAdministerClass failNode={<NotFoundPage />}>
      {children}
    </CanAdministerClass>
  );
}
