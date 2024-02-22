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
    title: `Class Group Manager: ${params.id}`,
    description: "OAMS Class Group Manager",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function AdminPanelClassGroupManagerLayout({
  children,
}: {
  children: ReactNode;
}) {
  return <>{children}</>;
}
