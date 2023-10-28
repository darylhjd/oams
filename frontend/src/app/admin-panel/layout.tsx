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
  return <>{children}</>;
}
