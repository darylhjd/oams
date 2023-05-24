import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:frontend/widgets/dialogs.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

class ProfileScreen extends ConsumerWidget {
  static const double mobilePadding = 10;
  static const double desktopPadding = 20;

  const ProfileScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    var userAsync = ref.watch(userInfoProvider);

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

  Widget mobile(BuildContext context, User data) {
    return ListView(
      padding: const EdgeInsets.all(mobilePadding),
      children: [
        _PreferredUsername(true, data.preferredUsername),
      ],
    );
  }

  Widget desktop(BuildContext context, User data) {
    return ListView(
      padding: const EdgeInsets.all(desktopPadding),
      children: [
        _PreferredUsername(false, data.preferredUsername),
      ],
    );
  }
}

class _PreferredUsername extends StatelessWidget {
  static const double mobileTopPadding = 50;
  static const double mobileOtherPadding = 5;
  static const double desktopTopPadding = 100;
  static const double desktopOtherPadding = 10;
  final bool isMobile;
  final String preferredUsername;

  const _PreferredUsername(this.isMobile, this.preferredUsername, {Key? key})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      alignment: Alignment.bottomCenter,
      padding: isMobile
          ? const EdgeInsets.fromLTRB(
              mobileOtherPadding,
              mobileTopPadding,
              mobileOtherPadding,
              mobileOtherPadding,
            )
          : const EdgeInsets.fromLTRB(
              desktopOtherPadding,
              desktopTopPadding,
              desktopOtherPadding,
              desktopOtherPadding,
            ),
      child: Text(
        preferredUsername,
        style: Theme.of(context)
            .textTheme
            .bodyLarge
            ?.copyWith(fontWeight: FontWeight.bold),
      ),
    );
  }
}
