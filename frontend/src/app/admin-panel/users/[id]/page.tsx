"use client";

import { User } from "@/api/user";
import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";

export default function AdminPanelUserPage({ params }: { params: Params }) {
  const [user, setUser] = useState<User>({} as User);
  const promiseFunc = async () => {
    const data = await APIClient.userGet(params.id);
    return setUser(data.user);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <UserDisplay user={user} />
    </RequestLoader>
  );
}

function UserDisplay({ user }: { user: User }) {
  return <Text ta="center">{user.id}</Text>;
}
