import 'dart:convert';

import 'package:app/network/interfaces.dart';
import 'package:app/screens/bulk-edit/model.dart';
import 'package:app/screens/bulk-edit/parse.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

const String _originalText =
    '''% Enter flashcards in CSV format, deliminated by the | char.
% The questions will be formatted using markdown, and the \\ escapes chars.
% Comments start with a `%`, \\n denotes a newline, and tables can be written with \\|
% For example, a question can be written as:
% | _Question 1_ | *Answer 1*\\n*Answer 2* | Tag-a Tag-b Tag-c
% In the above question, Question 1 will be italic, Answer 1 and Answer 2 will be bold on two seperate lines.
% Once you submit the question, the server will give it an ID, so the next time you see it, it will look like:
% 12345| _Question 1_ | *Answer 1*\\n*Answer 2* | Tag-a Tag-b Tag-c
% Questions and IDs must be unique per user\n''';

Future<String> getAllQuestions({required FetchData remoteServer}) async {
  final List<dynamic> dynamicFlashcards = jsonDecode(
        await remoteServer.getAllQuestions(),
      )['flashcards'] ??
      [];

  final flashcards =
      dynamicFlashcards.map((e) => Flashcard.fromJson(e)).toIList();
  return '$_originalText${unparse(flashcards)}';
}

Future<void> postAllQuestions(
    {required FetchData remoteServer,
    required IList<Flashcard> flashcards}) async {
  await remoteServer.postAllQuestions(flashcards);
}
