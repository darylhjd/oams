import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// Shows the profile screen.
class ProfileScreen extends ConsumerWidget {
  static const double _mobilePadding = 30;
  static const double _desktopHorizontalPadding = 20;
  static const double _desktopVerticalPadding = 40;

  const ProfileScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final userMeResponse = ref.watch(sessionUserProvider).requireValue;
    return ScreenTemplate(ResponsiveBreakpoints.of(context).isMobile
        ? _mobile(context, userMeResponse)
        : _desktop(context, userMeResponse));
  }

  Widget _mobile(BuildContext context, UserMeResponse data) {
    return ListView(
      padding: const EdgeInsets.symmetric(vertical: _mobilePadding),
      children: [
        const _ProfileHeader(),
        const SizedBox(height: _mobilePadding),
        _DetailsCard(data.sessionUser),
      ],
    );
  }

  Widget _desktop(BuildContext context, UserMeResponse data) {
    return Container(
      padding: const EdgeInsets.symmetric(
        horizontal: _desktopHorizontalPadding,
        vertical: _desktopVerticalPadding,
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            padding:
                const EdgeInsets.fromLTRB(0, 0, _desktopVerticalPadding, 0),
            child: const _ProfileHeader(),
          ),
          Flexible(
            child: ListView(
              children: [
                _DetailsCard(data.sessionUser),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// Shows the ID and name of the user.
class _ProfileHeader extends StatelessWidget {
  const _ProfileHeader();

  @override
  Widget build(BuildContext context) {
    return Text(
      "Your Profile",
      style: Theme.of(context)
          .textTheme
          .headlineLarge
          ?.copyWith(fontWeight: FontWeight.bold),
      textAlign: TextAlign.center,
    );
  }
}

// _DetailsCard shows the details of the user.
class _DetailsCard extends StatelessWidget {
  static const double _padding = 50;

  final User _user;

  const _DetailsCard(this._user);

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Container(
        padding: const EdgeInsets.all(_padding),
        child: Column(
          children: [
            TextFormField(
              decoration: const InputDecoration(labelText: "ID"),
              readOnly: true,
              initialValue: _user.id,
            ),
            TextFormField(
              decoration: const InputDecoration(labelText: "Name"),
              readOnly: true,
              initialValue: _user.name.isEmpty ? "-" : _user.name,
            ),
            TextFormField(
              decoration: const InputDecoration(labelText: "Email"),
              readOnly: true,
              initialValue: _user.email.isEmpty ? "-" : _user.email,
            ),
            TextFormField(
              decoration: const InputDecoration(labelText: "Role"),
              readOnly: true,
              initialValue: _user.role.name,
            ),
          ],
        ),
      ),
    );
  }
}
