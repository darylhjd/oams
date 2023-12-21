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
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
