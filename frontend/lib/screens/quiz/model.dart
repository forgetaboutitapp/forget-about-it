import '../../network/interfaces.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../fn/fn.dart';
import '../../protobufs-build/client_to_server.pb.dart' as client_to_server;
import '../../protobufs-build/server_to_client.pbenum.dart'
    as server_to_client_enums;

part 'model.g.dart';

part 'model.freezed.dart';

enum QuestionType { dueQuestion, nonDueQuestion, newQuestion }

@freezed
sealed class QuizQuestionState with _$QuizQuestionState {
  factory QuizQuestionState.waiting() = QuizQuestionStateWaiting;

  factory QuizQuestionState.none() = QuizQuestionStateNone;

  factory QuizQuestionState.data({
    required int id,
    required String question,
    required String answer,
    required QuestionType questionType,
    required int newCards,
    required int dueCards,
    required int nonDueCards,
  }) = QuizQuestionStateData;
}

@riverpod
class QuizQuestions extends _$QuizQuestions {
  Result<QuizQuestionState> _state = Ok(QuizQuestionState.waiting());

  @override
  Result<QuizQuestionState> build() {
    return _state;
  }

  Future<void> gradeQuestion(
    FetchDataWithToken client,
    ISet<String>? tagsQuery,
    int questionID,
    bool correct,
  ) async {
    _state = Ok(QuizQuestionState.waiting());
    ref.invalidateSelf();

    final p = await client.gradeQuestion(
      client_to_server.GradeQuestion(questionid: questionID, correct: correct),
    );
    (await p.doMap((_) async => await getNextQuestion(client, tagsQuery)))
        .match(
            onOk: (_) {},
            onErr: (e) {
              _state = Err(e);
              ref.invalidateSelf();
            });
  }

  Future<Result<()>> getNextQuestion(
    FetchDataWithToken client,
    ISet<String>? tagsQuery,
  ) async {
    if (tagsQuery == null) {
      _state = Err(NoTagException());
      ref.invalidateSelf();
      return Err(NoTagException());
    }
    _state = Ok(QuizQuestionState.waiting());
    ref.invalidateSelf();

    _state = switch (await client.getNextQuestion(
      client_to_server.GetNextQuestion(
        tags: tagsQuery.toIList(),
      ),
    )) {
      Err(value: final e) => Err(e),
      Ok(value: final v) => Ok(
          QuizQuestionStateData(
            id: v.flashcard.id,
            question: v.flashcard.question,
            answer: v.flashcard.answer,
            questionType: switch (v.typeOfQuestion) {
              server_to_client_enums
                    .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_NEW =>
                QuestionType.newQuestion,
              server_to_client_enums
                    .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_DUE =>
                QuestionType.dueQuestion,
              server_to_client_enums
                    .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_NON_DUE =>
                QuestionType.nonDueQuestion,
              server_to_client_enums.GetNextQuestion_TypeOfQuestion
                    .TYPE_OF_QUESTION_UNSPECIFIED =>
                throw UnimplementedError(),
              server_to_client_enums.GetNextQuestion_TypeOfQuestion() =>
                throw UnimplementedError(),
            },
            newCards: v.newQuestions,
            dueCards: v.dueQuestions,
            nonDueCards: v.nonDueQuestions,
          ),
        ),
    };
    ref.invalidateSelf();
    return Ok(());
  }
}

class NoTagException implements Exception {}
