import { Text } from "@mantine/core";

export type Params = {
  id: string;
};

export default function AdminPanelUsersPage({ params }: { params: Params }) {
  return <Text ta="center">{params.id}</Text>;
}
