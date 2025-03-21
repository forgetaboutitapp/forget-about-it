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
    try {
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
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<void> postAllQuestions(IList<Flashcard> flashcards) async {
    try {
      final res = await client.post(
        Uri.parse('$remoteHost/api/v0/secure/post-all-questions'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token'
        },
        body: jsonEncode(flashcards.map((e) => e.toJson()).toList()),
      );
      if (res.statusCode != 200) {
        throw ServerException(code: res.statusCode);
      }
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<String> generateNewToken() async {
    try {
      final v = await client.get(
          Uri.parse('$remoteHost/api/v0/secure/generate-new-token'),
          headers: {
            'Cache-Control': 'no-cache',
            'Authorization': 'Bearer $token'
          });
      return v.body;
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<bool> checkNewToken() async {
    try {
      final v = await client.get(
          Uri.parse('$remoteHost/api/v0/secure/check-new-token'),
          headers: {
            'Cache-Control': 'no-cache',
            'Authorization': 'Bearer $token'
          });
      if (v.body == 'done') {
        return true;
      }
      await Future.delayed(Duration(seconds: 1));
      return false;
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<void> deleteNewToken() async {
    try {
      await client.get(Uri.parse('$remoteHost/api/v0/secure/delete-new-token'),
          headers: {
            'Cache-Control': 'no-cache',
            'Authorization': 'Bearer $token'
          });
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<String> getRemoteSettings() async {
    try {
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
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  static Future<String> update(
      http.Client client, Uri remoteHost, Map<String, dynamic> d) async {
    try {
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
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  String getRemoteHost() => remoteHost;

  @override
  Future<String> getAllTags() async {
    try {
      final getAllTags = await client
          .get(Uri.parse('$remoteHost/api/v0/secure/get-all-tags'), headers: {
        'Cache-Control': 'no-cache',
        'Authorization': 'Bearer $token'
      });
      if (getAllTags.statusCode != 200) {
        throw ServerException(code: getAllTags.statusCode);
      }
      return getAllTags.body;
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<String> getNextQuestion(ISet<String> tags) async {
    try {
      final nextQuestion = await client.post(
        Uri.parse('$remoteHost/api/v0/secure/get-next-question'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode(tags.toList()),
      );
      if (nextQuestion.statusCode != 200) {
        throw ServerException(code: nextQuestion.statusCode);
      }
      return nextQuestion.body;
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<void> gradeQuestion(int questionID, bool correct) async {
    try {
      final nextQuestion = await client.post(
        Uri.parse('$remoteHost/api/v0/secure/grade-question'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode({'question-id': questionID, 'correct': correct}),
      );
      if (nextQuestion.statusCode != 200) {
        throw ServerException(code: nextQuestion.statusCode);
      }
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }

  @override
  Future<void> uploadAlgorithm(String data) async {
    try {
      final nextQuestion = await client.post(
        Uri.parse('$remoteHost/api/v0/secure/upload-algorithm'),
        headers: {
          'Cache-Control': 'no-cache',
          'Authorization': 'Bearer $token',
        },
        body: jsonEncode({'data': data}),
      );
      if (nextQuestion.statusCode != 200) {
        throw ServerException(code: nextQuestion.statusCode);
      }
    } on http.ClientException catch (_) {
      throw ServerException(code: -1);
    }
  }
}
