"use client";

import { Text } from "@mantine/core";
import { useState } from "react";
import { Params } from "./layout";
import { APIClient } from "@/api/client";
import { RequestLoader } from "@/components/request_loader";
import { Class } from "@/api/class";

export default function AdminPanelClassPage({ params }: { params: Params }) {
  const [cls, setClass] = useState<Class>({} as Class);
  const promiseFunc = async () => {
    const data = await APIClient.classGet(params.id);
    return setClass(data.class);
  };

  return (
    <RequestLoader promiseFunc={promiseFunc}>
      <ClassDisplay cls={cls} />
    </RequestLoader>
  );
}

function ClassDisplay({ cls }: { cls: Class }) {
  return <Text ta="center">{cls.id}</Text>;
}
