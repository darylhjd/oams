// LoginResponse is a data class modelling the response from a login API request.
class LoginResponse {
  static const String redirectUrlField = "redirect_url";

  final String redirectUrl;

  LoginResponse(this.redirectUrl);

  LoginResponse.fromJson(Map<String, dynamic> json)
      : redirectUrl = json[redirectUrlField];
}

// User models a response returned from the users API endpoint.
class SessionUserInfoResponse {
  final User? sessionUser;

  SessionUserInfoResponse.fromJson(Map<String, dynamic> json)
      : sessionUser = json["session_user"] == null
            ? null
            : User.fromJson(json["session_user"]);
}

// User models a user entity from the API.
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
