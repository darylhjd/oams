"use client";

import { User } from "@/api/user";
import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";

export default function AdminPanelUsersPage({ params }: { params: Params }) {
  const [user, setUser] = useState<User | null>(null);

  if (user == null) {
    return <Text>Loading or does not exist.</Text>;
  }

  return <Page user={user} />;
}

function Page({ user }: { user: User }) {
  return <Text ta="center">{user.id}</Text>;
}
