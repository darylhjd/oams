"use client";

import { Params } from "@/app/admin-panel/class-group-managers/[id]/layout";
import { Text } from "@mantine/core";
import { useState } from "react";
import { ClassGroupManager } from "@/api/class_group_manager";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";

export default function AdminPanelClassGroupManagerPage({
  params,
}: {
  params: Params;
}) {
  const [manager, setManager] = useState<ClassGroupManager>(
    {} as ClassGroupManager,
  );
  const promiseFunc = async () => {
    const data = await APIClient.classGroupManagerGet(params.id);
    return setManager(data.manager);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <ClassGroupManagerDisplay manager={manager} />
    </RequestLoader>
  );
}

function ClassGroupManagerDisplay({ manager }: { manager: ClassGroupManager }) {
  return <Text ta="center">{manager.id}</Text>;
}
