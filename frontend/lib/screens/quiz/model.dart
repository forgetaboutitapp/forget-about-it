import 'dart:convert';

import 'package:app/data/errors.dart';
import 'package:app/network/interfaces.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'model.g.dart';
part 'model.freezed.dart';

@freezed
sealed class QuizQuestionState with _$QuizQuestionState {
  factory QuizQuestionState.waiting() = QuizQuestionStateWaiting;
  factory QuizQuestionState.none() = QuizQuestionStateNone;
  factory QuizQuestionState.error({required Exception exception}) =
      QuizQuestionStateError;
  factory QuizQuestionState.data({
    required int id,
    required String question,
    required String answer,
  }) = QuizQuestionStateData;
}

QuizQuestionStateData _quizQuestionStateDatafromJson(
        Map<String, dynamic> data) =>
    QuizQuestionStateData(
        id: data['id'], question: data['question'], answer: data['answer']);

@riverpod
class QuizQuestions extends _$QuizQuestions {
  QuizQuestionState _state = QuizQuestionState.waiting();

  @override
  QuizQuestionState build() {
    return _state;
  }

  Future<void> gradeQuestion(
    FetchData client,
    ISet<String>? tagsQuery,
    int questionID,
    bool correct,
  ) async {
    try {
      _state = QuizQuestionState.waiting();
      ref.invalidateSelf();

      await client.gradeQuestion(questionID, correct);

      getNextQuestion(client, tagsQuery);
    } on ServerException catch (e) {
      _state = QuizQuestionState.error(exception: e);
      ref.invalidateSelf();
    }
  }

  Future<void> getNextQuestion(
    FetchData client,
    ISet<String>? tagsQuery,
  ) async {
    try {
      if (tagsQuery == null) {
        _state = QuizQuestionState.error(exception: NoTagException());
        ref.invalidateSelf();

        return;
      }
      _state = QuizQuestionState.waiting();
      ref.invalidateSelf();

      _state = _quizQuestionStateDatafromJson(
        jsonDecode(
          await client.getNextQuestion(tagsQuery),
        ),
      );
      ref.invalidateSelf();
    } on ServerException catch (e) {
      _state = QuizQuestionState.error(exception: e);
      ref.invalidateSelf();
    }
  }
}

class NoTagException implements Exception {}
