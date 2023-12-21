"use client";

import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { ClassGroup } from "@/api/class_group";

export default function AdminPanelClassGroupPage({
  params,
}: {
  params: Params;
}) {
  const [classGroup, setClassGroup] = useState<ClassGroup>({} as ClassGroup);
  const promiseFunc = async () => {
    const data = await APIClient.classGroupGet(params.id);
    return setClassGroup(data.class_group);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <ClassGroupDisplay classGroup={classGroup} />
    </RequestLoader>
  );
}

function ClassGroupDisplay({ classGroup }: { classGroup: ClassGroup }) {
  return <Text ta="center">{classGroup.id}</Text>;
}
