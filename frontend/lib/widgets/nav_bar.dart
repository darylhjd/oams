import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:frontend/api/client.dart';
import 'package:frontend/providers/session.dart';
import 'package:frontend/router.dart';
import 'package:go_router/go_router.dart';
import 'package:responsive_framework/responsive_framework.dart';

// The navigation bar for the app.
class NavBar extends StatelessWidget {
  static const double _height = 90;
  static const double _margin = 10;
  static const double _padding = 5;
  static final List<BoxShadow>? _boxShadow = kElevationToShadow[8];
  static const double _borderRadius = 10;

  const NavBar({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: _height,
      child: Container(
        margin: const EdgeInsets.all(_margin),
        padding: const EdgeInsets.all(_padding),
        decoration: BoxDecoration(
          color: Theme.of(context).primaryColorLight,
          boxShadow: _boxShadow,
          borderRadius: const BorderRadius.all(
            Radius.circular(_borderRadius),
          ),
        ),
        child: ResponsiveBreakpoints.of(context).isMobile
            ? _mobile(context)
            : _desktop(context),
      ),
    );
  }

  Widget _mobile(BuildContext context) {
    return const Row(
      mainAxisSize: MainAxisSize.max,
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        _Options(true),
        _Logo(),
      ],
    );
  }

  Widget _desktop(BuildContext context) {
    return const Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        _Logo(),
        _Options(false),
      ],
    );
  }
}

// Shows the app logo on the navigation bar. Doubles as a button to go back to
// the index screen on press.
class _Logo extends StatelessWidget {
  static const double _width = 100;
  static const double _height = double.infinity;

  static const String _logoPath = "assets/logo.png";

  const _Logo({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: _width,
      height: _height,
      child: IconButton(
        icon: Image.asset(_logoPath),
        onPressed: () => context.goNamed(Routes.index.name),
      ),
    );
  }
}

// Stores all the navigation bar items except the logo.
class _Options extends ConsumerWidget {
  final bool _isMobile;

  const _Options(this._isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final isLoggedIn = ref.read(sessionProvider);
    return _isMobile
        ? _mobile(context, isLoggedIn)
        : _desktop(context, isLoggedIn);
  }

  Widget _mobile(BuildContext context, bool isLoggedIn) {
    var items = <PopupMenuItem<Widget>>[
      const PopupMenuItem(
        padding: EdgeInsets.zero,
        child: _AboutButton(true),
      ),
    ];

    if (isLoggedIn) {
      items.addAll([
        const PopupMenuItem(
          padding: EdgeInsets.zero,
          child: _ProfileButton(true),
        ),
        const PopupMenuItem(
          padding: EdgeInsets.zero,
          child: _LogoutButton(),
        )
      ]);
    } else {
      items.add(const PopupMenuItem(
        padding: EdgeInsets.zero,
        child: _LoginButton(true),
      ));
    }

    return PopupMenuButton(
      icon: const Icon(Icons.menu),
      itemBuilder: (context) => items,
    );
  }

  Widget _desktop(BuildContext context, bool isLoggedIn) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        const _AboutButton(false),
        isLoggedIn ? const _ProfileButton(false) : const _LoginButton(false),
      ],
    );
  }
}

// Shows the About button on the navigation bar.
class _AboutButton extends StatelessWidget {
  final Widget _text = const _NavBarText("About");
  final bool _isMobile;

  const _AboutButton(this._isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return PopupMenuItem<Widget>(
      onTap: () => context.goNamed(Routes.about.name),
      child: _text,
    );
  }

  Widget _desktop(BuildContext context) {
    return TextButton(
      onPressed: () => context.goNamed(Routes.about.name),
      child: _text,
    );
  }
}

// Shows the Profile button on the navigation bar.
class _ProfileButton extends StatelessWidget {
  final Widget _text = const _NavBarText("Profile");
  final bool _isMobile;

  const _ProfileButton(this._isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return PopupMenuItem<Widget>(
      child: _text,
      onTap: () => context.goNamed(Routes.profile.name),
    );
  }

  Widget _desktop(BuildContext context) {
    return PopupMenuButton(
      icon: const Icon(Icons.account_circle_rounded),
      itemBuilder: (context) => [
        PopupMenuItem(
          onTap: () => context.goNamed(Routes.profile.name),
          child: _text,
        ),
        const PopupMenuItem(
          padding: EdgeInsets.zero,
          child: _LogoutButton(),
        ),
      ],
    );
  }
}

// Shows the Login button on the navigation bar.
class _LoginButton extends StatelessWidget {
  final Widget _text = const _NavBarText("Login");
  final bool _isMobile;

  const _LoginButton(this._isMobile, {Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return _isMobile ? _mobile(context) : _desktop(context);
  }

  Widget _mobile(BuildContext context) {
    return PopupMenuItem<Widget>(
      child: _text,
      onTap: () => context.goNamed(Routes.login.name),
    );
  }

  Widget _desktop(BuildContext context) {
    return TextButton(
      child: _text,
      onPressed: () => context.goNamed(Routes.login.name),
    );
  }
}

// Shows the Logout button on the navigation bar.
class _LogoutButton extends ConsumerWidget {
  final Widget _text = const _NavBarText("Logout", color: Colors.red);

  const _LogoutButton({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return PopupMenuItem(
      onTap: () async {
        await APIClient.logout();
        ref.read(sessionProvider.notifier).update((_) => false);
        ref.invalidate(sessionUserProvider);

        // Do not depend on context mounting for routing in async function.
        ref.read(routerProvider).goNamed(Routes.index.name);
      },
      child: Text.rich(
        TextSpan(
          children: [
            const WidgetSpan(
              child: Icon(
                Icons.logout,
                color: Colors.red,
              ),
            ),
            WidgetSpan(
              child: _text,
            ),
          ],
        ),
      ),
    );
  }
}

// Helps to standardise the text style in the navigation bar.
class _NavBarText extends StatelessWidget {
  final String _text;
  final Color? _color;

  const _NavBarText(this._text, {Key? key, Color? color})
      : _color = color,
        super(key: key);

  @override
  Widget build(BuildContext context) {
    return Text(
      _text,
      style: Theme.of(context).textTheme.titleMedium?.copyWith(color: _color),
    );
  }
}
