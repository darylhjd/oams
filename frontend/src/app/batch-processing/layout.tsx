import { UserRole } from "@/api/user";
import { CheckHasUserRole } from "@/components/session_checker";
import { ReactNode } from "react";

export const metadata = {
  title: "Batch Processing",
  description: "OAMS Batch Processing Service",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function BatchProcessingLayout({
  children,
}: {
  children: ReactNode;
}) {
  return (
    <CheckHasUserRole role={UserRole.SystemAdmin}>{children}</CheckHasUserRole>
  );
}
