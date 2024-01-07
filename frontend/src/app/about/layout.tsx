import { ReactNode } from "react";

export const metadata = {
  title: "About",
  description: "About OAMS",
  icons: {
    icon: "/favicon.svg",
  },
};

export default function AboutLayout({ children }: { children: ReactNode }) {
  return <>{children}</>;
}
