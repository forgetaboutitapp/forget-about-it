import 'dart:developer' as developer;

import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';

class FutureWidget<T> extends HookWidget {
  final Future<T> future;
  final Widget Function(BuildContext context, T r) built;
  final Widget Function(BuildContext context) waiting;
  const FutureWidget({
    super.key,
    required this.future,
    required this.built,
    required this.waiting,
  });

  @override
  Widget build(BuildContext context) {
    final ValueNotifier<T?> result = useState(null);
    useEffect(() {
      () async {
        final T v = await future;
        result.value = v;
      }();
      return () {};
    });
    final r = result.value;
    if (r == null) {
      return waiting(context);
    } else {
      developer.log('finished building $r');
      return built(context, r);
    }
  }
}
