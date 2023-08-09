import 'dart:convert';
import 'dart:io';

import 'package:frontend/api/models.dart';
import 'package:frontend/env/env.dart';
import 'package:http/browser_client.dart';

class APIClient {
  static final _client = () {
    var client = BrowserClient();
    client.withCredentials = true;
    return client;
  }();

  static final Uri _apiUri = Uri.parse("${apiServerHost()}:${apiServerPort()}");
  static final String _defaultRedirectUrl =
      "${webServerHost()}:${webServerPort()}";
  static const String loginRedirectUrlParam = "redirect_url";

  static const String _loginPath = "api/v1/login";
  static const String _logoutPath = "api/v1/logout";
  static const String _userMePath = "api/v1/users/me";
  static const String _userPath = "api/v1/users/";

  // Get login URL to redirect user to SSO login site.
  static Future<String> getLoginURL(String returnTo) async {
    if (returnTo.isEmpty) {
      returnTo = _defaultRedirectUrl;
    }

    final uri = _apiUri.replace(
      path: _loginPath,
      queryParameters: {
        loginRedirectUrlParam: returnTo,
      },
    );

    final response = await _client.get(uri);
    final body = jsonDecode(response.body);

    if (response.statusCode != HttpStatus.ok) {
      return Future.error(HttpException(ErrorResponse.fromJson(body).message));
    }

    return LoginResponse.fromJson(body).redirectUrl;
  }

  // Remove the current user session, and also helps unset the session cookie.
  static Future<bool> logout() async {
    final uri = _apiUri.replace(
      path: _logoutPath,
    );

    final response = await _client.get(uri);
    return response.statusCode == HttpStatus.ok;
  }

  // Get current user session information. For more information on the user, use
  // the GET endpoint for user.
  static Future<UserMeResponse> getUserMe() async {
    final uri = _apiUri.replace(
      path: _userMePath,
    );

    final response = await _client.get(uri);
    final body = jsonDecode(response.body);

    if (response.statusCode != HttpStatus.ok) {
      return Future.error(HttpException(ErrorResponse.fromJson(body).message));
    }

    return UserMeResponse.fromJson(body);
  }

  // Get a user information by ID.
  static Future<GetUserResponse> getUser(String id) async {
    final uri = _apiUri.replace(
      path: "$_userPath$id",
    );

    final response = await (_client.get(uri));
    final body = jsonDecode(response.body);

    if (response.statusCode != HttpStatus.ok) {
      return Future.error(HttpException(ErrorResponse.fromJson(body).message));
    }

    return GetUserResponse.fromJson(body);
  }
}
