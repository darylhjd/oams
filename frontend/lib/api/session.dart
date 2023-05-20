import 'package:flutter_riverpod/flutter_riverpod.dart';

final sessionProvider = FutureProvider<User>((ref) async {
  return User();
});

class User {}
