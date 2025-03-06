import 'dart:convert';

import 'package:app/data/errors.dart';
import 'package:app/screens/settings/model.dart';
import 'package:http/http.dart' as http;

Future<RemoteSettings> getRemoteSettings(
    String remoteHost, String token, http.Client client) async {
  final remoteSettings = await client.get(
      Uri.parse('$remoteHost/api/v0/secure/get-remote-settings'),
      headers: {'Cache-Control': 'no-cache', 'Authorization': 'Bearer $token'});
  if (remoteSettings.statusCode != 200) {
    throw ServerException(code: remoteSettings.statusCode);
  }
  final settings = jsonDecode(remoteSettings.body);
  return RemoteSettings.fromJSON(settings);
}
