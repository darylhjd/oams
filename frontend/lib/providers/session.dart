import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/api/models.dart';

// sessionUserProvider provides information about the session
// by getting the current user info.
// We do this since the API uses a HttpOnly cookie to store
// the required auth codes and so we cannot read it directly.
final sessionUserProvider = FutureProvider<User>((ref) async {
  try {
    var resp = await APIClient.getUsers();
    if (resp.sessionUser == null) {
      throw Exception();
    }
    ref.read(sessionProvider.notifier).update((_) => true);
    return resp.sessionUser!;
  } catch (e) {
    ref.read(sessionProvider.notifier).update((_) => false);
    rethrow;
  }
});

// sessionProvider provides a simple boolean flag to check if
// the current session has a user logged in or not.
final sessionProvider = StateProvider<bool>((ref) {
  return false;
});
