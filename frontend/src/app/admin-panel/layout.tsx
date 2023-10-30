import { UserRole } from "@/api/user";
import { CheckHasUserRole } from "@/components/session_checker";

export const metadata = {
  title: "Admin Panel",
  description: "OAMS Admin Panel",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function AdminPanelLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <CheckHasUserRole role={UserRole.SystemAdmin}>{children}</CheckHasUserRole>
  );
}
