import 'dart:convert';

import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/foundation.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:http/http.dart' as http;
import '../main.dart';

part 'login.g.dart';
part 'login.freezed.dart';

sealed class LoginReturn {}

class NoLogin extends LoginReturn {}

@freezed
class LoggedIn extends LoginReturn with _$LoggedIn {
  const factory LoggedIn({required String id}) = _LoggedIn;
}

@Riverpod(keepAlive: true)
class Login extends _$Login {
  @override
  LoginReturn build() {
    String? str = sharedPreferences.getString('LOGGED_IN_KEY');
    if (str == null) {
      return NoLogin();
    }
    try {
      return LoggedIn(id: str);
    } catch (e) {
      // This means that the uuid is corrupted
      sharedPreferences.remove('LOGGED_IN_KEY');
      return NoLogin();
    }
  }

  Future<bool> update(
      http.Client client, Uri remoteHost, IList<String> twelveWords) async {
    final v = await client.post(
      Uri(
        scheme: remoteHost.scheme,
        port: remoteHost.port,
        host: remoteHost.host,
        path: '/api/v0/get-token/by-twelve-words',
      ),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: jsonEncode({'twelve-words': twelveWords.toList()}),
    );
    Map<String, dynamic> json = jsonDecode(v.body);
    if (!json.containsKey('token')) {
      return false;
    }
    await sharedPreferences.setString('LOGGED_IN_KEY', json['token']);
    ref.invalidateSelf();
    return true;
  }
}
