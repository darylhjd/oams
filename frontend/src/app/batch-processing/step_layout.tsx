import { Divider, Space } from "@mantine/core";
import React from "react";

export function StepLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Divider my="sm" />
      <Space visibleFrom="md" h="md" />
      {children}
    </>
  );
}
