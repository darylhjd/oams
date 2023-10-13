import { APIClient } from "@/api/client";
import { AsyncDataSource } from "@/components/entity_tables";

export class UsersDataSource extends AsyncDataSource {
  async getRows(offset: number, limit: number): Promise<any[]> {
    const response = await APIClient.usersGet(offset, limit);
    if (response == null) {
      return [];
    }

    super.updateRecordsEstimationState(offset, limit, response.users.length);
    return response.users;
  }
}

export class ClassesDataSource extends AsyncDataSource {
  async getRows(offset: number, limit: number): Promise<any[]> {
    const response = await APIClient.classesGet(offset, limit);
    if (response == null) {
      return [];
    }

    super.updateRecordsEstimationState(offset, limit, response.classes.length);
    return response.classes;
  }
}

export class ClassGroupsDataSource extends AsyncDataSource {
  async getRows(offset: number, limit: number): Promise<any[]> {
    const response = await APIClient.classGroupsGet(offset, limit);
    if (response == null) {
      return [];
    }

    super.updateRecordsEstimationState(
      offset,
      limit,
      response.class_groups.length,
    );
    return response.class_groups;
  }
}

export class PlaceholderDataSource extends AsyncDataSource {
  async getRows(offset: number, limit: number): Promise<any[]> {
    return [];
  }
}
