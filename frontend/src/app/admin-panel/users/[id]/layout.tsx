import { UserRole } from "@/api/user";
import { CheckHasUserRole } from "@/components/session_checker";
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
  return (
    <CheckHasUserRole role={UserRole.SystemAdmin}>{children}</CheckHasUserRole>
  );
}
