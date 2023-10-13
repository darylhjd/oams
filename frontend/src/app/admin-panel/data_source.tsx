import { APIClient } from "@/api/client";
import { AsyncDataSource } from "@/components/entity_tables";

export class UsersDataSource extends AsyncDataSource {
  async getRows(offset: number, limit: number): Promise<any[]> {
    console.log("users data source called!");
    const response = await APIClient.usersGet(offset, limit);
    if (response == null) {
      return [];
    }

    super.updateRecordsEstimationState(offset, limit, response.users.length);
    return response.users;
  }
}

export class PlaceholderDataSource extends AsyncDataSource {
  async getRows(offset: number, limit: number): Promise<any[]> {
    return [];
  }
}
