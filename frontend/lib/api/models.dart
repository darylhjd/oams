// LoginResponse is a data class modelling the response from a login API request.
class LoginResponse {
  static const String redirectUrlField = "redirect_url";

  final String redirectUrl;

  LoginResponse(this.redirectUrl);

  LoginResponse.fromJson(Map<String, dynamic> json)
      : redirectUrl = json[redirectUrlField];
}

// User models a user returned from the API.
class User {
  final String homeAccountID;
  final String preferredUsername;

  User.fromJson(Map<String, dynamic> json)
      : homeAccountID = json["home_account_id"],
        preferredUsername = json["username"];
}
