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

  User.fromJson(Map<String, dynamic> json)
      : id = json["id"],
        name = json["name"],
        email = json["email"],
        role = UserRole.fromValue(json["role"]);
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
