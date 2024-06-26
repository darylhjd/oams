import {
  MRT_TableOptions,
  MRT_PaginationState,
  MantineReactTable,
  MRT_DensityState,
  useMantineReactTable,
} from "mantine-react-table";
import {
  AsyncDataSource,
  ClassGroupSessionsDataSource,
  ClassGroupsDataSource,
  ClassGroupManagersDataSource,
  ClassesDataSource,
  SessionEnrollmentsDataSource,
  UsersDataSource,
  ClassAttendanceRulesDataSource,
} from "./data_sources";
import { useEffect, useState } from "react";
import {
  ClassesDataTableColumns,
  ClassGroupsDataTableColumns,
  ClassGroupSessionsDataTableColumns,
  ClassGroupManagersDataTableColumns,
  SessionEnrollmentsDataTableColumns,
  UsersDataTableColumns,
  DEFAULT_PAGE_SIZE,
  ClassAttendanceRulesDataTableColumns,
} from "@/components/columns";
import { useRouter } from "next/navigation";
import { Routes } from "@/routing/routes";

const ROW_PROPS = {
  style: { cursor: "pointer" },
};

export function UsersTable() {
  const router = useRouter();

  return (
    <AsyncDataTable
      columns={UsersDataTableColumns}
      dataSource={new UsersDataSource()}
      mantineTableBodyRowProps={({ row }) => ({
        onClick: (_) => router.push(Routes.adminPanelUser + row.original.id),
        ...ROW_PROPS,
      })}
    />
  );
}

export function ClassesTable() {
  const router = useRouter();

  return (
    <AsyncDataTable
      columns={ClassesDataTableColumns}
      dataSource={new ClassesDataSource()}
      mantineTableBodyRowProps={({ row }) => ({
        onClick: (_) => router.push(Routes.adminPanelClass + row.original.id),
        ...ROW_PROPS,
      })}
    />
  );
}

export function ClassAttendanceRulesTable() {
  const router = useRouter();

  return (
    <AsyncDataTable
      columns={ClassAttendanceRulesDataTableColumns}
      dataSource={new ClassAttendanceRulesDataSource()}
      mantineTableBodyRowProps={({ row }) => ({
        onClick: (_) =>
          router.push(Routes.adminPanelClassAttendanceRule + row.original.id),
        ...ROW_PROPS,
      })}
    />
  );
}

export function ClassGroupsTable() {
  const router = useRouter();

  return (
    <AsyncDataTable
      columns={ClassGroupsDataTableColumns}
      dataSource={new ClassGroupsDataSource()}
      mantineTableBodyRowProps={({ row }) => ({
        onClick: (_) =>
          router.push(Routes.adminPanelClassGroup + row.original.id),
        ...ROW_PROPS,
      })}
    />
  );
}

export function ClassGroupManagersTable() {
  const router = useRouter();

  return (
    <AsyncDataTable
      columns={ClassGroupManagersDataTableColumns}
      dataSource={new ClassGroupManagersDataSource()}
      mantineTableBodyRowProps={({ row }) => ({
        onClick: (_) =>
          router.push(Routes.adminPanelClassGroupManager + row.original.id),
        ...ROW_PROPS,
      })}
    />
  );
}

export function ClassGroupSessionsTable() {
  const router = useRouter();

  return (
    <AsyncDataTable
      columns={ClassGroupSessionsDataTableColumns}
      dataSource={new ClassGroupSessionsDataSource()}
      mantineTableBodyRowProps={({ row }) => ({
        onClick: (_) =>
          router.push(Routes.adminPanelClassGroupSession + row.original.id),
        ...ROW_PROPS,
      })}
    />
  );
}

export function SessionEnrollmentsTable() {
  const router = useRouter();

  return (
    <AsyncDataTable
      columns={SessionEnrollmentsDataTableColumns}
      dataSource={new SessionEnrollmentsDataSource()}
      mantineTableBodyRowProps={({ row }) => ({
        onClick: (_) =>
          router.push(Routes.adminPanelSessionEnrollment + row.original.id),
        ...ROW_PROPS,
      })}
    />
  );
}

function AsyncDataTable<T extends Record<string, any>>({
  dataSource,
  ...props
}: {
  dataSource: AsyncDataSource<T>;
} & Omit<MRT_TableOptions<T>, "data">) {
  const [data, setData] = useState<T[]>([]);
  const [loading, setLoading] = useState(true);

  const [rowCount, setRowCount] = useState(0);
  const [paginationState, setPaginationState] = useState<MRT_PaginationState>({
    pageIndex: 0,
    pageSize: DEFAULT_PAGE_SIZE,
  });

  useEffect(() => {
    setLoading(true);

    const from = paginationState.pageIndex * paginationState.pageSize;
    dataSource.getRows(from, paginationState.pageSize).then((data) => {
      setData(data);
      setRowCount(dataSource.totalRecords);
      setLoading(false);
    });
  }, [paginationState]);

  const table = useMantineReactTable({
    ...props,
    data,
    manualPagination: true,
    rowCount: rowCount,
    initialState: {
      density: "sm" as MRT_DensityState,
    },
    state: {
      pagination: paginationState,
      isLoading: loading,
    },
    onPaginationChange: setPaginationState,
    enableStickyHeader: true,
    enableStickyFooter: true,
  });

  return <MantineReactTable table={table} />;
}
