"use client";

import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { EntityLoader } from "@/components/entity_loader";
import { ClassGroupSession } from "@/api/class_group_session";

export default function AdminPanelClassGroupSessionPage({
  params,
}: {
  params: Params;
}) {
  const [classGroupSession, setClassGroupSession] =
    useState<ClassGroupSession | null>(null);
  const promiseFunc = async () => {
    const data = await APIClient.classGroupSessionGet(params.id);
    return setClassGroupSession(data.class_group_session);
  };

  return (
    <EntityLoader promiseFunc={promiseFunc}>
      <ClassGroupSessionDisplay classGroupSession={classGroupSession!} />
    </EntityLoader>
  );
}

function ClassGroupSessionDisplay({
  classGroupSession,
}: {
  classGroupSession: ClassGroupSession;
}) {
  return <Text ta="center">{classGroupSession.id}</Text>;
}
