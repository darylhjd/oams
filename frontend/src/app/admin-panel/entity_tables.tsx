import {
  AsyncDataTable,
  ClassDataTableColumns,
  ClassGroupDataTableColumns,
  ClassGroupSessionDataTableColumns,
  SessionEnrollmentsDataTableColumns,
  UserDataTableColumns,
} from "@/components/entity_tables";
import { PlaceholderDataSource, UsersDataSource } from "./data_source";

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
      dataSource={new PlaceholderDataSource()}
    />
  );
}

export function ClassGroupsTable() {
  return (
    <AsyncDataTable
      columns={ClassGroupDataTableColumns}
      dataSource={new PlaceholderDataSource()}
    />
  );
}

export function ClassGroupSessionsTable() {
  return (
    <AsyncDataTable
      columns={ClassGroupSessionDataTableColumns}
      dataSource={new PlaceholderDataSource()}
    />
  );
}

export function SessionEnrollmentsTable() {
  return (
    <AsyncDataTable
      columns={SessionEnrollmentsDataTableColumns}
      dataSource={new PlaceholderDataSource()}
    />
  );
}
