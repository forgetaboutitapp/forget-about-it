import 'package:web/web.dart' as web;

String? getCurrentLocation() =>
    '${web.window.location.protocol}//${web.window.location.hostname}:${web.window.location.port}';
