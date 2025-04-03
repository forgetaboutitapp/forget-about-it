import 'package:flutter_web_plugins/flutter_web_plugins.dart';
import 'package:web/web.dart' as web;

String? getCurrentLocation() => switch (web.window.location.port) {
      '' => '${web.window.location.protocol}//${web.window.location.hostname}',
      _ =>
        '${web.window.location.protocol}//${web.window.location.hostname}:${web.window.location.port}'
    };

void usePathUrlStrategySafe() => usePathUrlStrategy();
