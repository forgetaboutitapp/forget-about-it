import 'dart:convert';

import 'package:app/screens/login/submit_type.dart';
import 'package:flutter/foundation.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:hive_ce/hive.dart';
import 'package:http/http.dart' as http;
import '../data/constants.dart';
import '../network/network.dart';

part 'login.freezed.dart';

@freezed
sealed class LoginReturn with _$LoginReturn {
  const factory LoginReturn.noLogin() = _LoginReturnNoLogin;
  const factory LoginReturn.loggedIn({required String token}) =
      _LoginReturnLoggedIn;
}

Future<bool> update(
    http.Client client, Uri remoteUri, SubmitType submitType) async {
  Map<String, dynamic> d = switch (submitType) {
    TwelveWords(:final twelveWords) => {'twelve-words': twelveWords.toList()},
    Token(:final token) => {'token': token},
  };
  final body = await RemoteServer.update(client, remoteUri, d);
  Map<String, dynamic> json = jsonDecode(body);
  if (!json.containsKey('token')) {
    return false;
  }
  Hive.box(localSettingsHiveBox)
      .put(localSettingsHiveLoginToken, json['token']);
  Hive.box(localSettingsHiveBox).put(localSettingsHiveRemoteHost,
      '${remoteUri.scheme}://${remoteUri.host}:${remoteUri.port}');
  return true;
}
