import 'package:forget_about_it/protobufs-build/client_server/v1/client_to_server.pbgrpc.dart';

import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../../fn/fn.dart';
import '../../interop/grpc_channel.dart';
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
            createGrpcChannel(Uri.parse(host)))
        .gradeQuestion(GradeQuestionRequest(
            token: token, questionid: questionID, correct: correct));
    if (client.hasError()) {
      if (client.error.shouldLogOut) {
        logOut();
      }
      _state = Err(Exception(client.error.error));
      ref.invalidateSelf();
      return null;
    }

    if (!client.hasOk()) {
      _state = Err(Exception('Server Error'));
      ref.invalidateSelf();
      return null;
    }

    final nextDue = client.ok.nextDue.toInt();
    final nextQuestion = await getNextQuestion(
      host,
      token,
      logOut,
      tagsQuery,
      forceNewQuestion,
    );
    return switch (nextQuestion) {
      Ok() => nextDue,
      Err(:final value) => (() {
          _state = Err(value);
          ref.invalidateSelf();
          return null;
        })(),
    };
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
            createGrpcChannel(Uri.parse(host)))
        .getNextQuestion(GetNextQuestionRequest(
            token: token,
            tags: tagsQuery.toIList(),
            getNewQuestion: forceNewQuestion));
    if (client.hasError()) {
      if (client.error.shouldLogOut) {
        logOut();
      }
      final error = Exception(client.error.error);
      _state = Err(error);
      ref.invalidateSelf();
      return Err(error);
    }

    if (!client.hasOk()) {
      final error = Exception('Server Error');
      _state = Err(error);
      ref.invalidateSelf();
      return Err(error);
    }

    final question = client.ok;
    _state = Ok(
      QuizQuestionStateData(
        id: question.flashcard.id,
        question: question.flashcard.question,
        answer: question.flashcard.answer,
        questionType: switch (question.typeOfQuestion) {
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_NEW =>
            QuestionType.newQuestion,
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_DUE =>
            QuestionType.dueQuestion,
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_NON_DUE =>
            QuestionType.nonDueQuestion,
          server_to_client_enums
                .GetNextQuestion_TypeOfQuestion.TYPE_OF_QUESTION_UNSPECIFIED =>
            throw UnimplementedError(),
          server_to_client_enums.GetNextQuestion_TypeOfQuestion() =>
            throw UnimplementedError(),
        },
        newCards: question.newQuestions,
        dueCards: question.dueQuestions,
        nonDueCards: question.nonDueQuestions,
      ),
    );
    ref.invalidateSelf();
    return Ok(());
  }
}

class NoTagException implements Exception {}
