import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/api/models.dart';

// Provides information about the current session use. We do this since the API
// uses a HttpOnly cookie to store the required auth codes and so we cannot read
// it directly.
final sessionUserProvider = FutureProvider<UserMeResponse>((ref) async {
  try {
    var resp = await APIClient.getUserMe();
    ref.read(sessionProvider.notifier).update((_) => true);
    return resp;
  } catch (e) {
    ref.read(sessionProvider.notifier).update((_) => false);
    rethrow;
  }
});

// Provides a simple boolean flag to check if the current session has a user
// logged in or not.
final sessionProvider = StateProvider<bool>((ref) => false);
