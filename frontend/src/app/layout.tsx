import type { Metadata } from 'next'
import React, { Suspense } from "react";
import Providers from "@/app/providers";
import Header from '@/components/header';

export const metadata: Metadata = {
  title: 'OAMS',
  description: 'Online Attendance Management System',
  icons: {
    icon: '/favicon.svg',
  },
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body>
        <Providers>
          <Header />
          {children}
        </Providers>
      </body>
    </html>
  )
}
