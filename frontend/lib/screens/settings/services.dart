import 'dart:convert';

import '../../network/interfaces.dart';
import 'models/remote_settings.dart';

Future<RemoteSettings> getRemoteSettings(FetchData remoteServer) async =>
    RemoteSettings.fromJSON(jsonDecode(await remoteServer.getRemoteSettings()));
