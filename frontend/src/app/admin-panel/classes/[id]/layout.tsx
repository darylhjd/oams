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
    title: `Class: ${params.id}`,
    description: "OAMS Class",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function AdminPanelClassLayout({
  children,
}: {
  children: ReactNode;
}) {
  return <>{children}</>;
}
