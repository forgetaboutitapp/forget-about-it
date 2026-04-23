import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';
import 'package:forget_about_it/protobufs-build/client_server/v1/server_to_client.pb.dart';
import 'package:grpc/grpc_web.dart';

import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../fn/fn.dart';
import '../../protobufs-build/client_server/v1/server_to_client.pbenum.dart'
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

  Future<int?> gradeQuestion(
    String host,
    String token,
    Function logOut,
    ISet<String>? tagsQuery,
    int questionID,
    bool correct,
    bool forceNewQuestion,
  ) async {
    _state = Ok(QuizQuestionState.waiting());
    ref.invalidateSelf();
    final client = await ForgetAboutItServiceClient(
            GrpcWebClientChannel.xhr(Uri.parse(host)))
        .gradeQuestion(GradeQuestionRequest(
            token: token, questionid: questionID, correct: correct));
    int? whenNextQuestionDue;
    final _ = switch (client) {
      GradeQuestion(:final nextDue) => (() async {
          final _ = switch (await getNextQuestion(
            host,
            token,
            logOut,
            tagsQuery,
            forceNewQuestion,
          )) {
            Ok() => (() {
                whenNextQuestionDue = nextDue.toInt();
                return null;
              })(),
            Err(:final value) => (() {
                _state = Err(Exception(value));
                ref.invalidateSelf();
              })()
          };
        })(),
      Err(:final value) => (() {
          _state = Err(value);
          ref.invalidateSelf();
          return null;
        })(),
      var v => (() {
          _state = Err(Exception('State Error: $v'));
          ref.invalidateSelf();
          return null;
        })(),
    };
    int? returnVal = whenNextQuestionDue;
    return returnVal;
  }

  Future<Result<()>> getNextQuestion(
    String host,
    String token,
    Function logOut,
    ISet<String>? tagsQuery,
    bool forceNewQuestion,
  ) async {
    if (tagsQuery == null) {
      _state = Err(NoTagException());
      ref.invalidateSelf();
      return Err(NoTagException());
    }
    _state = Ok(QuizQuestionState.waiting());
    ref.invalidateSelf();

    final client = await ForgetAboutItServiceClient(
            GrpcWebClientChannel.xhr(Uri.parse(host)))
        .getNextQuestion(GetNextQuestionRequest(
            token: token,
            tags: tagsQuery.toIList(),
            getNewQuestion: forceNewQuestion));
    _state = switch (client) {
      Err(:final value) => Err(Exception(value)),
      GetNextQuestion(
        :final newQuestions,
        :final dueQuestions,
        :final nonDueQuestions,
        :final flashcard,
        :final typeOfQuestion,
      ) =>
        Ok(
          QuizQuestionStateData(
            id: flashcard.id,
            question: flashcard.question,
            answer: flashcard.answer,
            questionType: switch (typeOfQuestion) {
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
            newCards: newQuestions,
            dueCards: dueQuestions,
            nonDueCards: nonDueQuestions,
          ),
        ),
      var v => Err(Exception('State Error: $v')),
    };
    ref.invalidateSelf();
    return Ok(());
  }
}

class NoTagException implements Exception {}
