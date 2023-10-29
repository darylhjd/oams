import { IsLoggedIn } from "@/components/session_checker";

export const metadata = {
  title: "Profile",
  description: "Current Session User",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function ProfileLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <IsLoggedIn>{children}</IsLoggedIn>;
}
