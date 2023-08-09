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
  final User? sessionUser;

  UserMeResponse.fromJson(Map<String, dynamic> json)
      : sessionUser = json["session_user"] == null
            ? null
            : User.fromJson(json["session_user"]);
}

class User {
  final String id;
  final String name;
  final String email;
  final String role;

  User.fromJson(Map<String, dynamic> json)
      : id = json["id"],
        name = json["name"],
        email = json["email"],
        role = json["role"];
}

class GetUserResponse {
  final User user;

  GetUserResponse.fromJson(Map<String, dynamic> json)
      : user = User.fromJson(json["user"]);
}
