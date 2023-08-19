import 'dart:io';

import 'package:dio/browser.dart';
import 'package:dio/dio.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/env/env.dart';

class APIClient {
  static final _client = Dio()
    ..httpClientAdapter = BrowserHttpClientAdapter(withCredentials: true);

  static final Uri _apiUri = Uri.parse("${apiServerHost()}:${apiServerPort()}");
  static final String _defaultRedirectUrl =
      "${webServerHost()}:${webServerPort()}";
  static const String loginRedirectUrlParam = "redirect_url";

  static const String _loginPath = "api/v1/login";
  static const String _logoutPath = "api/v1/logout";

  static const String _usersPath = "api/v1/users";
  static const String _userPath = "api/v1/users/";
  static const String _userMePath = "api/v1/users/me";

  static const String _classesPath = "api/v1/classes";

  static const String _classGroupsPath = "api/v1/class-groups";

  static const String _classGroupSessionsPath = "api/v1/class-group-sessions";

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

    final response = await _client.getUri(uri);
    return LoginResponse.fromJson(response.data).redirectUrl;
  }

  // Remove the current user session, and also helps unset the session cookie.
  static Future<bool> logout() async {
    final uri = _apiUri.replace(
      path: _logoutPath,
    );

    final response = await _client.getUri(uri);
    return response.statusCode == HttpStatus.ok;
  }

  // Get current user session information. For more information on the user, use
  // the GET endpoint for user.
  static Future<UserMeResponse> getUserMe() async {
    final uri = _apiUri.replace(
      path: _userMePath,
    );

    final response = await _client.getUri(uri);
    return UserMeResponse.fromJson(response.data);
  }

  // Get a user information by ID.
  static Future<GetUserResponse> getUser(String id) async {
    final uri = _apiUri.replace(
      path: "$_userPath$id",
    );

    final response = await _client.getUri(uri);
    return GetUserResponse.fromJson(response.data);
  }

  // Get users.
  static Future<GetUsersResponse> getUsers(int limit, int offset) async {
    final uri = _apiUri.replace(
      path: _usersPath,
      queryParameters: {
        "limit": limit.toString(),
        "offset": offset.toString(),
      },
    );

    final response = await _client.getUri(uri);
    return GetUsersResponse.fromJson(response.data);
  }

  static Future<GetClassesResponse> getClasses(int limit, int offset) async {
    final uri = _apiUri.replace(
      path: _classesPath,
      queryParameters: {
        "limit": limit.toString(),
        "offset": offset.toString(),
      },
    );

    final response = await _client.getUri(uri);
    return GetClassesResponse.fromJson(response.data);
  }

  static Future<GetClassGroupsResponse> getClassGroups(
      int limit, int offset) async {
    final uri = _apiUri.replace(
      path: _classGroupsPath,
      queryParameters: {
        "limit": limit.toString(),
        "offset": offset.toString(),
      },
    );

    final response = await _client.getUri(uri);
    return GetClassGroupsResponse.fromJson(response.data);
  }

  static Future<GetClassGroupSessionsResponse> getClassGroupSessions(
      int limit, int offset) async {
    final uri =
        _apiUri.replace(path: _classGroupSessionsPath, queryParameters: {
      "limit": limit.toString(),
      "offset": offset.toString(),
    });

    final response = await _client.getUri(uri);
    return GetClassGroupSessionsResponse.fromJson(response.data);
  }
}
