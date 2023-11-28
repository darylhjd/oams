"use client";

import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { EntityLoader } from "@/app/admin-panel/entity_loader";
import { ClassGroup } from "@/api/class_group";

export default function AdminPanelClassGroupPage({
  params,
}: {
  params: Params;
}) {
  const [classGroup, setClassGroup] = useState<ClassGroup | null>(null);
  const promiseFunc = async () => {
    const data = await APIClient.classGroupGet(params.id);
    return setClassGroup(data.class_group);
  };

  return (
    <EntityLoader promiseFunc={promiseFunc}>
      <ClassGroupDisplay classGroup={classGroup!} />
    </EntityLoader>
  );
}

function ClassGroupDisplay({ classGroup }: { classGroup: ClassGroup }) {
  return <Text ta="center">{classGroup.id}</Text>;
}
