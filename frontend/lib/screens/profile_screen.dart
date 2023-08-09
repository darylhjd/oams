import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:frontend/widgets/dialogs.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// Shows the profile screen.
class ProfileScreen extends ConsumerWidget {
  static const double _mobilePadding = 30;
  static const double _desktopPadding = 40;

  const ProfileScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return ScreenTemplate(
      ref.watch(sessionUserProvider).when(
            data: (data) => ResponsiveBreakpoints.of(context).isMobile
                ? _mobile(context, data)
                : _desktop(context, data),
            loading: () => const Center(
              child: CircularProgressIndicator(),
            ),
            error: (error, stackTrace) => const InvalidSession(),
          ),
    );
  }

  Widget _mobile(BuildContext context, User user) {
    return ListView(
      padding: const EdgeInsets.symmetric(vertical: _mobilePadding),
      children: [
        _NameHeader(true, user),
        const SizedBox(height: _mobilePadding),
        _DetailsCard(true, user),
      ],
    );
  }

  Widget _desktop(BuildContext context, User user) {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: _desktopPadding),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            padding: const EdgeInsets.fromLTRB(0, 0, _desktopPadding, 0),
            child: _NameHeader(false, user),
          ),
          Flexible(
            child: ListView(
              children: [
                _DetailsCard(false, user),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

// Shows the ID and name of the user.
class _NameHeader extends StatelessWidget {
  final User _user;
  final bool _isMobile;

  const _NameHeader(this._isMobile, this._user);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return Column(
      children: [
        Text(
          _user.id,
          style: Theme.of(context)
              .textTheme
              .headlineLarge
              ?.copyWith(fontWeight: FontWeight.bold),
        ),
        Text(
          _user.name,
          style: Theme.of(context).textTheme.bodyLarge,
        ),
      ],
    );
  }

  Widget _desktop(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          _user.id,
          style: Theme.of(context)
              .textTheme
              .headlineLarge
              ?.copyWith(fontWeight: FontWeight.bold),
        ),
        Text(
          _user.name,
          style: Theme.of(context).textTheme.bodyLarge,
        ),
      ],
    );
  }
}

// _DetailsCard shows the details of the user.
class _DetailsCard extends StatelessWidget {
  final User _user;
  final bool _isMobile;

  const _DetailsCard(this._isMobile, this._user);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return const Placeholder();
  }

  Widget desktop(BuildContext context) {
    return const Placeholder();
  }
}
