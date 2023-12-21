import { Metadata } from "next";
import React from "react";

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
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
