import {
  AsyncDataTable,
  ClassDataTableColumns,
  ClassGroupDataTableColumns,
  ClassGroupSessionDataTableColumns,
  SessionEnrollmentsDataTableColumns,
  UserDataTableColumns,
} from "@/components/entity_tables";

export function UsersTable() {
  return <AsyncDataTable columns={UserDataTableColumns} />;
}

export function ClassesTable() {
  return <AsyncDataTable columns={ClassDataTableColumns} />;
}

export function ClassGroupsTable() {
  return <AsyncDataTable columns={ClassGroupDataTableColumns} />;
}

export function ClassGroupSessionsTable() {
  return <AsyncDataTable columns={ClassGroupSessionDataTableColumns} />;
}

export function SessionEnrollmentsTable() {
  return <AsyncDataTable columns={SessionEnrollmentsDataTableColumns} />;
}
