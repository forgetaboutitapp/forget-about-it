import 'dart:convert';

import 'package:app/data/errors.dart';
import 'package:app/screens/bulk-edit/model.dart';
import 'package:app/screens/bulk-edit/parse.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:http/http.dart' as http;

const String _originalText =
    '''% Enter flashcards in CSV format, deliminated by the | char.
% The questions will be formatted using markdown, and the \\ escapes chars.
% Comments start with a `%`, \\n denotes a newline, and tables can be written with \\|
% For example, a question can be written as:
% | _Question 1_ | *Answer 1*\\n*Answer 2* | Tag-a Tag-b Tag-c|
% In the above question, Question 1 will be italic, Answer 1 and Answer 2 will be bold on two seperate lines.
% Once you submit the question, the server will give it an ID, so the next time you see it, it will look like:
% fcb54af1-2c6f-474d-ad95-0c1b170a2991| _Question 1_ | *Answer 1*\\n*Answer 2* | Tag-a Tag-b Tag-c|
% Questions and IDs must be unique per user\n''';

Future<String> getAllQuestions(
    http.Client client, String token, String remoteHost) async {
  final res = await client.get(
      Uri.parse('$remoteHost/api/v0/secure/get-all-questions'),
      headers: {'Cache-Control': 'no-cache', 'Authorization': 'Bearer $token'});
  if (res.statusCode != 200) {
    throw ServerException(code: res.statusCode);
  }
  final List<dynamic> dynamicFlashcards = jsonDecode(res.body);

  final flashcards =
      dynamicFlashcards.map((e) => Flashcard.fromJson(e)).toIList();
  return '$_originalText${unparse(flashcards)}';
}

Future<void> postAllQuestions(http.Client client, String token,
    String remoteHost, IList<Flashcard> flashcards) async {
  print('data: ${jsonEncode(flashcards.map((e) => e.toJson()).toList())}');
  final res = await client.post(
    Uri.parse('$remoteHost/api/v0/secure/post-all-questions'),
    headers: {'Cache-Control': 'no-cache', 'Authorization': 'Bearer $token'},
    body: jsonEncode(flashcards.map((e) => e.toJson()).toList()),
  );
  if (res.statusCode != 200) {
    throw ServerException(code: res.statusCode);
  }
}
