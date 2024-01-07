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
    title: `Class Group: ${params.id}`,
    description: "OAMS Class Group",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function AdminPanelClassGroupLayout({
  children,
}: {
  children: ReactNode;
}) {
  return <>{children}</>;
}
