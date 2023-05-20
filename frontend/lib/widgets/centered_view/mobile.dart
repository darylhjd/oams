import 'package:flutter/material.dart';

class CenteredViewMobile extends StatelessWidget {
  final Widget child;

  const CenteredViewMobile({Key? key, required this.child}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 5),
      alignment: Alignment.center,
      child: child,
    );
  }
}
