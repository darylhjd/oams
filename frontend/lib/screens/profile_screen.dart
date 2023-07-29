import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:frontend/widgets/dialogs.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

// ProfileScreen shows the profile screen.
class ProfileScreen extends ConsumerWidget {
  static const double mobilePadding = 30;
  static const double desktopPadding = 20;

  const ProfileScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    var userAsync = ref.watch(sessionUserProvider);

    return ScreenTemplate(
      userAsync.when(
        data: (data) => ResponsiveBreakpoints.of(context).isMobile
            ? mobile(context, data)
            : desktop(context, data),
        loading: () => const Center(
          child: CircularProgressIndicator(),
        ),
        error: (error, stackTrace) => const InvalidSession(),
      ),
    );
  }

  Widget mobile(BuildContext context, User user) {
    return ListView(
      padding: const EdgeInsets.symmetric(vertical: mobilePadding),
      children: [
        _NameHeader(true, user),
      ],
    );
  }

  Widget desktop(BuildContext context, User user) {
    return ListView(
      padding: const EdgeInsets.all(desktopPadding),
      children: [
        _NameHeader(false, user),
      ],
    );
  }
}

// _NameHeader shows as a larger text the name of a User.
class _NameHeader extends StatelessWidget {
  final User user;
  final bool isMobile;

  const _NameHeader(this.isMobile, this.user);

  @override
  Widget build(BuildContext context) {
    return isMobile ? mobile(context) : desktop(context);
  }

  Widget mobile(BuildContext context) {
    return Column(
      children: [
        Text(
          user.id,
          style: Theme.of(context)
              .textTheme
              .headlineLarge
              ?.copyWith(fontWeight: FontWeight.bold),
        ),
        Text(
          user.name,
          style: Theme.of(context).textTheme.bodyLarge,
        ),
      ],
    );
  }

  Widget desktop(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          user.id,
          style: Theme.of(context)
              .textTheme
              .headlineLarge
              ?.copyWith(fontWeight: FontWeight.bold),
        ),
        Text(
          user.name,
          style: Theme.of(context).textTheme.bodyLarge,
        ),
      ],
    );
  }
}
