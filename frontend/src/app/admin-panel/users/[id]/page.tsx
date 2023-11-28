"use client";

import { User } from "@/api/user";
import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { EntityLoader } from "@/app/admin-panel/entity_loader";

export default function AdminPanelUserPage({ params }: { params: Params }) {
  const [user, setUser] = useState<User | null>(null);
  const promiseFunc = async () => {
    const data = await APIClient.userGet(params.id);
    return setUser(data.user);
  };

  return (
    <EntityLoader promiseFunc={promiseFunc}>
      <UserDisplay user={user!} />
    </EntityLoader>
  );
}

function UserDisplay({ user }: { user: User }) {
  return <Text ta="center">{user.id}</Text>;
}
