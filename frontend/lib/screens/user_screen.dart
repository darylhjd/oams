import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:frontend/widgets/dialogs.dart';

final _userProvider = FutureProvider.autoDispose
    .family<GetUserResponse, String>((ref, userId) async {
  return await APIClient.getUser(userId);
});

// Shows information about a user.
class UserScreen extends ConsumerWidget {
  final String _userId;

  const UserScreen(this._userId, {super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ScreenTemplate(ref.watch(_userProvider(_userId)).when(
          data: (data) => _screen(context, data),
          error: (error, stackTrace) => const InvalidSession(),
          loading: () => const Center(child: CircularProgressIndicator()),
        ));
  }

  Widget _screen(BuildContext context, GetUserResponse data) {
    return ListView(children: [
      _UserHeader(data),
    ]);
  }
}

class _UserHeader extends StatelessWidget {
  final GetUserResponse _data;

  const _UserHeader(this._data);

  @override
  Widget build(BuildContext context) {
    return Column(children: [
      Text(_data.user.id),
    ]);
  }
}
