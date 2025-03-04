import 'dart:convert';

import 'package:app/data/errors.dart';
import 'package:app/screens/login/submit_type.dart';
import 'package:flutter/foundation.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:hive_ce/hive.dart';
import 'package:http/http.dart' as http;
import '../data/constants.dart';

part 'login.freezed.dart';

sealed class LoginReturn {}

class NoLogin extends LoginReturn {}

@freezed
class LoggedIn extends LoginReturn with _$LoggedIn {
  const factory LoggedIn({required String id}) = _LoggedIn;
}

Future<bool> update(
    http.Client client, Uri remoteHost, SubmitType submitType) async {
  Map<String, dynamic> d = switch (submitType) {
    TwelveWords(:final twelveWords) => {'twelve-words': twelveWords.toList()},
    Token(:final token) => {'token': token},
  };
  final host = Uri(
    scheme: remoteHost.scheme,
    port: remoteHost.port,
    host: remoteHost.host,
  );
  final v = await client.post(
    Uri(
      scheme: remoteHost.scheme,
      port: remoteHost.port,
      host: remoteHost.host,
      path: '/api/v0/get-token',
    ),
    headers: <String, String>{
      'Content-Type': 'application/json; charset=UTF-8',
    },
    body: jsonEncode(d),
  );
  if (v.statusCode != 200) {
    throw ServerException(code: v.statusCode);
  }
  Map<String, dynamic> json = jsonDecode(v.body);
  if (!json.containsKey('token')) {
    return false;
  }
  Hive.box(localSettingsHiveBox)
      .put(localSettingsHiveLoginToken, json['token']);
  Hive.box(localSettingsHiveBox)
      .put(localSettingsHiveRemoteHost, host.toString());
  return true;
}
