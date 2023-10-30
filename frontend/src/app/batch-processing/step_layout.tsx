import { Divider, Space } from "@mantine/core";

export function StepLayout({ children }: { children: React.ReactNode }) {
  return (
    <>
      <Divider my="sm" />
      <Space visibleFrom="md" h="md" />
      {children}
    </>
  );
}
