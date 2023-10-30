import { Class } from "@/api/class";
import { ClassGroup } from "@/api/class_group";
import { ClassGroupSession } from "@/api/class_group_session";
import { ClassManager } from "@/api/class_manager";
import { APIClient } from "@/api/client";
import { SessionEnrollment } from "@/api/session_enrollment";
import { User } from "@/api/user";

abstract class AsyncDataSource<T> {
  constructor(
    public totalRecords: number = 0,
    public isApproximateRowCount: boolean = true,
  ) {}

  updateRecordsEstimationState(
    offset: number,
    limit: number,
    lastFetchLength: number,
  ) {
    const knownLength = offset + lastFetchLength;

    if (knownLength < offset + limit) {
      this.totalRecords = knownLength;
    } else {
      this.totalRecords = Math.max(this.totalRecords, offset + 2 * limit); // Allows possible fetch of next page.
    }
  }

  abstract getRows(offset: number, limit: number): Promise<T[]>;
}

export class UsersDataSource extends AsyncDataSource<User> {
  async getRows(offset: number, limit: number): Promise<User[]> {
    const response = await APIClient.usersGet(offset, limit);
    super.updateRecordsEstimationState(offset, limit, response.users.length);
    return response.users;
  }
}

export class ClassesDataSource extends AsyncDataSource<Class> {
  async getRows(offset: number, limit: number): Promise<Class[]> {
    const response = await APIClient.classesGet(offset, limit);
    super.updateRecordsEstimationState(offset, limit, response.classes.length);
    return response.classes;
  }
}

export class ClassManagersDataSource extends AsyncDataSource<ClassManager> {
  async getRows(offset: number, limit: number): Promise<ClassManager[]> {
    const response = await APIClient.classManagersGet(offset, limit);
    super.updateRecordsEstimationState(
      offset,
      limit,
      response.class_managers.length,
    );
    return response.class_managers;
  }
}

export class ClassGroupsDataSource extends AsyncDataSource<ClassGroup> {
  async getRows(offset: number, limit: number): Promise<ClassGroup[]> {
    const response = await APIClient.classGroupsGet(offset, limit);
    super.updateRecordsEstimationState(
      offset,
      limit,
      response.class_groups.length,
    );
    return response.class_groups;
  }
}

export class ClassGroupSessionsDataSource extends AsyncDataSource<ClassGroupSession> {
  async getRows(offset: number, limit: number): Promise<ClassGroupSession[]> {
    const response = await APIClient.classGroupSessionsGet(offset, limit);
    super.updateRecordsEstimationState(
      offset,
      limit,
      response.class_group_sessions.length,
    );
    return response.class_group_sessions;
  }
}

export class SessionEnrollmentsDataSource extends AsyncDataSource<SessionEnrollment> {
  async getRows(offset: number, limit: number): Promise<SessionEnrollment[]> {
    const response = await APIClient.sessionEnrollmentsGet(offset, limit);
    super.updateRecordsEstimationState(
      offset,
      limit,
      response.session_enrollments.length,
    );
    return response.session_enrollments;
  }
}
