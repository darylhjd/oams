import { Metadata } from "next";

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

export default function AdminPanelUsersLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
