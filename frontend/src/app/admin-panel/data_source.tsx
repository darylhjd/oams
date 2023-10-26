import { APIClient } from "@/api/client";
import {
  Class,
  ClassGroup,
  ClassGroupSession,
  ClassManager,
  SessionEnrollment,
  User,
} from "@/api/models";
import { AsyncDataSource } from "@/components/entity_tables";

export class UsersDataSource extends AsyncDataSource<User> {
  async getRows(offset: number, limit: number): Promise<User[]> {
    const response = await APIClient.usersGet(offset, limit);
    if (response == null) {
      return [];
    }

    super.updateRecordsEstimationState(offset, limit, response.users.length);
    return response.users;
  }
}

export class ClassesDataSource extends AsyncDataSource<Class> {
  async getRows(offset: number, limit: number): Promise<Class[]> {
    const response = await APIClient.classesGet(offset, limit);
    if (response == null) {
      return [];
    }

    super.updateRecordsEstimationState(offset, limit, response.classes.length);
    return response.classes;
  }
}

export class ClassManagersDataSource extends AsyncDataSource<ClassManager> {
  async getRows(offset: number, limit: number): Promise<ClassManager[]> {
    const response = await APIClient.classManagersGet(offset, limit);
    if (response == null) {
      return [];
    }

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

export class ClassGroupSessionsDataSource extends AsyncDataSource<ClassGroupSession> {
  async getRows(offset: number, limit: number): Promise<ClassGroupSession[]> {
    const response = await APIClient.classGroupSessionsGet(offset, limit);
    if (response == null) {
      return [];
    }

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
    if (response == null) {
      return [];
    }

    super.updateRecordsEstimationState(
      offset,
      limit,
      response.session_enrollments.length,
    );
    return response.session_enrollments;
  }
}
