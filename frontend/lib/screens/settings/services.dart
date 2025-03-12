import 'dart:convert';

import 'package:app/screens/settings/model.dart';

import '../../network/interfaces.dart';

Future<RemoteSettings> getRemoteSettings(FetchData remoteServer) async =>
    RemoteSettings.fromJSON(jsonDecode(await remoteServer.getRemoteSettings()));
