import 'dart:convert';

import 'package:app/data/errors.dart';
import 'package:app/network/interfaces.dart';
import 'package:app/screens/bulk-edit/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:http/http.dart' as http;

class RemoteServer implements FetchData {
  final String remoteHost;
  final String token;
  final http.Client client;
  RemoteServer({
    required this.remoteHost,
    required this.token,
    required this.client,
  });

  @override
  Future<String> getAllQuestions() async {
    final res = await client.get(
        Uri.parse('$remoteHost/api/v0/secure/get-all-questions'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token'
        });
    if (res.statusCode != 200) {
      throw ServerException(code: res.statusCode);
    }
    return res.body;
  }

  @override
  Future<void> postAllQuestions(IList<Flashcard> flashcards) async {
    final res = await client.post(
      Uri.parse('$remoteHost/api/v0/secure/post-all-questions'),
      headers: {'Cache-Control': 'no-cache', 'Authorization': 'Bearer $token'},
      body: jsonEncode(flashcards.map((e) => e.toJson()).toList()),
    );
    if (res.statusCode != 200) {
      throw ServerException(code: res.statusCode);
    }
  }

  @override
  Future<String> generateNewToken() async {
    final v = await client.get(
        Uri.parse('$remoteHost/api/v0/secure/generate-new-token'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token'
        });
    return v.body;
  }

  @override
  Future<bool> checkNewToken() async {
    final v = await client
        .get(Uri.parse('$remoteHost/api/v0/secure/check-new-token'), headers: {
      'Cache-Control': 'no-cache',
      'Authorization': 'Bearer $token'
    });
    if (v.body == 'done') {
      return true;
    }
    await Future.delayed(Duration(seconds: 1));
    return false;
  }

  @override
  Future<void> deleteNewToken() async {
    await client.get(Uri.parse('$remoteHost/api/v0/secure/delete-new-token'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token'
        });
  }

  @override
  Future<String> getRemoteSettings() async {
    final remoteSettings = await client.get(
        Uri.parse('$remoteHost/api/v0/secure/get-remote-settings'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token'
        });
    if (remoteSettings.statusCode != 200) {
      throw ServerException(code: remoteSettings.statusCode);
    }
    return remoteSettings.body;
  }

  static Future<String> update(
      http.Client client, Uri remoteHost, Map<String, dynamic> d) async {
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
    return v.body;
  }

  @override
  String getRemoteHost() => remoteHost;
}
