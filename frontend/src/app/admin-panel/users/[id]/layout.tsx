import { UserRole } from "@/api/user";
import { CheckHasUserRole } from "@/components/session_checker";
import { Metadata } from "next";
import { Params } from "./page";

export async function generateMetadata({
  params,
}: {
  params: Params;
}): Promise<Metadata> {
  return {
    title: params.id,
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
