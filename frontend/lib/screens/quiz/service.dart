import '../../network/interfaces.dart';
import '../../screens/quiz/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';

import '../../fn/fn.dart';
import '../../protobufs-build/client_to_server.pb.dart' as client_to_server;
import '../../protobufs-build/server_to_client.pbenum.dart'
    as server_to_client_enums;

Future<Result<QuizQuestionStateData>> getNextQuestion(
    FetchDataWithToken client, ISet<String> tagsQuery) async {
  return switch (await client.getNextQuestion(
      client_to_server.GetNextQuestion(tags: tagsQuery.toList()))) {
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
      )
  };
}
