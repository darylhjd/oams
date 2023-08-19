import 'package:frontend/api/client.dart';

class ErrorResponse {
  final String message;

  ErrorResponse.fromJson(Map<String, dynamic> json) : message = json["message"];
}

class LoginResponse {
  final String redirectUrl;

  LoginResponse.fromJson(Map<String, dynamic> json)
      : redirectUrl = json[APIClient.loginRedirectUrlParam];
}

class UserMeResponse {
  final User sessionUser;
  final List<UpcomingClassGroupSession> upcomingSessions;

  UserMeResponse.fromJson(Map<String, dynamic> json)
      : sessionUser = User.fromJson(json["session_user"]),
        upcomingSessions = List<UpcomingClassGroupSession>.from(
          (json["upcoming_class_group_sessions"] as List)
              .map((i) => UpcomingClassGroupSession.fromJson(i)),
        );
}

class User {
  final String id;
  final String name;
  final String email;
  final UserRole role;
  final DateTime createdAt;
  final DateTime updatedAt;

  User.fromJson(Map<String, dynamic> json)
      : id = json["id"],
        name = json["name"],
        email = json["email"],
        role = UserRole.fromValue(json["role"]),
        createdAt = DateTime.parse(json["created_at"]).toLocal(),
        updatedAt = DateTime.parse(json["updated_at"]).toLocal();
}

enum UserRole {
  student("STUDENT"),
  courseCoordinator("COURSE_COORDINATOR"),
  admin("ADMIN");

  final String name;

  const UserRole(this.name);

  factory UserRole.fromValue(String value) {
    return UserRole.values.firstWhere((e) => e.name == value);
  }
}

class UpcomingClassGroupSession {
  final String code;
  final int year;
  final String semester;
  final String name;
  final ClassType classType;
  final DateTime startTime;
  final DateTime endTime;

  UpcomingClassGroupSession.fromJson(Map<String, dynamic> json)
      : code = json["code"],
        year = json["year"],
        semester = json["semester"],
        name = json["name"],
        classType = ClassType.fromValue(json["class_type"]),
        startTime = DateTime.parse(json["start_time"]).toLocal(),
        endTime = DateTime.parse(json["end_time"]).toLocal();
}

enum ClassType {
  lec("LEC"),
  tut("TUT"),
  lab("LAB");

  final String name;

  const ClassType(this.name);

  factory ClassType.fromValue(String value) {
    return ClassType.values.firstWhere((e) => e.name == value);
  }
}

class GetUserResponse {
  final User user;

  GetUserResponse.fromJson(Map<String, dynamic> json)
      : user = User.fromJson(json["user"]);
}

class GetUsersResponse {
  final bool result;
  final List<User> users;

  GetUsersResponse.fromJson(Map<String, dynamic> json)
      : result = json["result"],
        users = List<User>.from(
            (json["users"] as List).map((i) => User.fromJson(i)));
}

class Class {
  final int id;
  final String code;
  final int year;
  final String semester;
  final String programme;
  final int au;
  final DateTime createdAt;
  final DateTime updatedAt;

  Class.fromJson(Map<String, dynamic> json)
      : id = json["id"],
        code = json["code"],
        year = json["year"],
        semester = json["semester"],
        programme = json["programme"],
        au = json["au"],
        createdAt = DateTime.parse(json["created_at"]).toLocal(),
        updatedAt = DateTime.parse(json["updated_at"]).toLocal();
}

class GetClassesResponse {
  final bool result;
  final List<Class> classes;

  GetClassesResponse.fromJson(Map<String, dynamic> json)
      : result = json["result"],
        classes = List<Class>.from(
            (json["classes"] as List).map((i) => Class.fromJson(i)));
}

class ClassGroup {
  final int id;
  final int classId;
  final String name;
  final ClassType classType;
  final DateTime createdAt;
  final DateTime updatedAt;

  ClassGroup.fromJson(Map<String, dynamic> json)
      : id = json["id"],
        classId = json["class_id"],
        name = json["name"],
        classType = ClassType.fromValue(json["class_type"]),
        createdAt = DateTime.parse(json["created_at"]).toLocal(),
        updatedAt = DateTime.parse(json["updated_at"]).toLocal();
}

class GetClassGroupsResponse {
  final bool result;
  final List<ClassGroup> classGroups;

  GetClassGroupsResponse.fromJson(Map<String, dynamic> json)
      : result = json["result"],
        classGroups = List<ClassGroup>.from(
            (json["class_groups"] as List).map((i) => ClassGroup.fromJson(i)));
}
