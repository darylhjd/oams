import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/models.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/screens/screen_template.dart';
import 'package:frontend/widgets/dialogs.dart';
import 'package:responsive_framework/responsive_breakpoints.dart';

class ProfileScreen extends ConsumerWidget {
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
      padding: const EdgeInsets.all(10),
      children: const [
        Text("mobile view!"),
      ],
    );
  }

  Widget desktop(BuildContext context, User data) {
    return ListView(
      padding: const EdgeInsets.all(20),
      children: const [
        Text("desktop view!"),
      ],
    );
  }
}
