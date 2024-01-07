import { Metadata } from "next";
import { ReactNode } from "react";
import { CanAdministerClass } from "@/components/session_checker";
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
    title: `Class Administration: ${params.id}`,
    description: "OAMS Class Administration",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function ClassAdministrationLayout({
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
