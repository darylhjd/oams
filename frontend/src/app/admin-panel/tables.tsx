import {
  MRT_ColumnDef,
  MRT_PaginationState,
  MantineReactTable,
} from "mantine-react-table";
import {
  AsyncDataSource,
  ClassGroupSessionsDataSource,
  ClassGroupsDataSource,
  ClassManagersDataSource,
  ClassesDataSource,
  SessionEnrollmentsDataSource,
  UsersDataSource,
} from "./data_sources";
import { useEffect, useState } from "react";
import {
  ClassDataTableColumns,
  ClassGroupDataTableColumns,
  ClassGroupSessionDataTableColumns,
  ClassManagerDataTableColumns,
  SessionEnrollmentsDataTableColumns,
  UserDataTableColumns,
} from "@/components/entity_columns";

const DEFAULT_PAGE_SIZE = 10;

export function UsersTable() {
  return (
    <AsyncDataTable
      columns={UserDataTableColumns}
      dataSource={new UsersDataSource()}
    />
  );
}

export function ClassesTable() {
  return (
    <AsyncDataTable
      columns={ClassDataTableColumns}
      dataSource={new ClassesDataSource()}
    />
  );
}

export function ClassManagersTable() {
  return (
    <AsyncDataTable
      columns={ClassManagerDataTableColumns}
      dataSource={new ClassManagersDataSource()}
    />
  );
}

export function ClassGroupsTable() {
  return (
    <AsyncDataTable
      columns={ClassGroupDataTableColumns}
      dataSource={new ClassGroupsDataSource()}
    />
  );
}

export function ClassGroupSessionsTable() {
  return (
    <AsyncDataTable
      columns={ClassGroupSessionDataTableColumns}
      dataSource={new ClassGroupSessionsDataSource()}
    />
  );
}

export function SessionEnrollmentsTable() {
  return (
    <AsyncDataTable
      columns={SessionEnrollmentsDataTableColumns}
      dataSource={new SessionEnrollmentsDataSource()}
    />
  );
}

function AsyncDataTable<T extends Record<string, any>>({
  columns,
  dataSource,
}: {
  columns: MRT_ColumnDef<T>[];
  dataSource: AsyncDataSource<T>;
}) {
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

  return (
    <MantineReactTable
      columns={columns}
      data={data}
      manualPagination={true}
      rowCount={rowCount}
      state={{
        pagination: paginationState,
        isLoading: loading,
      }}
      onPaginationChange={setPaginationState}
    />
  );
}
