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
    title: `Session Enrollment: ${params.id}`,
    description: "OAMS Session Enrollment",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function AdminPanelSessionEnrollmentLayout({
  children,
}: {
  children: ReactNode;
}) {
  return <>{children}</>;
}
