import { Metadata } from "next";
import { ReactNode } from "react";

export type Params = {
  id: string;
};

export async function generateMetadata({
  params,
}: {
  params: Params;
}): Promise<Metadata> {
  return {
    title: `User: ${params.id}`,
    description: "OAMS User",
    icons: {
      icon: "/favicon.svg",
    },
  };
}

export default function AdminPanelUserLayout({
  children,
}: {
  children: ReactNode;
}) {
  return <>{children}</>;
}
