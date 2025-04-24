import 'dart:convert';

import 'package:app/network/interfaces.dart';
import 'package:app/screens/quiz/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

Future<QuizQuestionStateData> getNextQuestion(
    FetchData client, ISet<String> tagsQuery) async {
  return quizQuestionStateDatafromJson(
    jsonDecode(
      await client.getNextQuestion(tagsQuery),
    ),
  );
}
