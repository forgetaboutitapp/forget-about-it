import 'package:app/data/errors.dart';
import 'package:app/network/interfaces.dart';
import 'package:app/screens/quiz/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/material.dart';
import 'package:flutter_hooks/flutter_hooks.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

import '../../fn/fn.dart';
import '../general-display/show_error.dart';

class QuizView extends HookConsumerWidget {
  static String location = '/quiz';
  final FetchDataWithToken remoteServer;
  final Map<String, List<String>> tags;
  final bool isDarkMode;
  const QuizView({
    super.key,
    required this.remoteServer,
    required this.tags,
    required this.isDarkMode,
  });
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final question = ref.watch(quizQuestionsProvider);
    final initialized = useState(false);
    final tagsSet = useMemoized(() {
      return tags['tags']?.toISet();
    });
    useEffect(() {
      if (initialized.value == false) {
        try {
          ref
              .read(quizQuestionsProvider.notifier)
              .getNextQuestion(remoteServer, tagsSet);
        } on ServerException catch (e) {
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                backgroundColor: Colors.red,
                content: Text('$e'),
              ),
            );
          }
        }
        initialized.value = true;
      }
      return () {};
    }(), [key]);
    return Scaffold(
      appBar: AppBar(
        title: Text('Quiz'),
      ),
      body: Center(
        child: switch (question) {
          Err(value: final exception) => ErrorScreen(exception: exception),
          Ok(value: final v) => switch (v) {
              QuizQuestionStateWaiting() => CircularProgressIndicator(),
              QuizQuestionStateNone() => Center(
                  child: Text('No questions available'),
                ),
              QuizQuestionStateData(
                :final question,
                :final answer,
                :final id,
                :final questionType,
                :final dueCards,
                :final nonDueCards,
                :final newCards,
              ) =>
                DisplayQuestion(
                  remoteServer: remoteServer,
                  question: question,
                  answer: answer,
                  id: id,
                  questionType: questionType,
                  tagsSet: tagsSet,
                  isDarkMode: isDarkMode,
                  amountDueQuestions: dueCards,
                  amountNewQuestions: newCards,
                  amountNonDueQuestions: nonDueCards,
                ),
            },
        },
      ),
    );
  }
}

class ErrorScreen extends StatelessWidget {
  const ErrorScreen({
    super.key,
    required this.exception,
  });

  final Exception exception;

  @override
  Widget build(BuildContext context) {
    WidgetsBinding.instance
        .addPostFrameCallback((_) => showError(context, exception.toString()));
    return Center(child: Text(''));
  }
}

class DisplayQuestion extends HookConsumerWidget {
  final String question;
  final String answer;
  final int id;
  final ISet<String>? tagsSet;
  final FetchDataWithToken remoteServer;
  final QuestionType questionType;
  final bool isDarkMode;
  final int amountNewQuestions;
  final int amountDueQuestions;
  final int amountNonDueQuestions;
  const DisplayQuestion({
    super.key,
    required this.question,
    required this.answer,
    required this.id,
    required this.tagsSet,
    required this.remoteServer,
    required this.questionType,
    required this.isDarkMode,
    required this.amountNewQuestions,
    required this.amountDueQuestions,
    required this.amountNonDueQuestions,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final shouldShowQuestion = useState(true);
    return Padding(
      padding: const EdgeInsets.all(1.0),
      child: Container(
        decoration: BoxDecoration(
          border: Border.all(
            color: switch (questionType) {
              QuestionType.dueQuestion =>
                isDarkMode ? Colors.red[200] : Colors.red,
              QuestionType.nonDueQuestion =>
                isDarkMode ? Colors.green[200] : Colors.green,
              QuestionType.newQuestion =>
                isDarkMode ? Colors.blue[200] : Colors.blue,
            }!, // Border color
            width: 4.0, // Border width
          ),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Center(
              child: Text(
                switch (questionType) {
                  QuestionType.dueQuestion => 'Due Question',
                  QuestionType.nonDueQuestion => 'Review Non Due Question',
                  QuestionType.newQuestion => 'New Question',
                },
              ),
            ),
            Expanded(
              child: Center(
                child: shouldShowQuestion.value
                    ? Markdown(
                        data: question,
                        selectable: true,
                      )
                    : Markdown(
                        data: answer,
                        selectable: true,
                      ),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    'New: $amountNewQuestions',
                    style: TextStyle(
                      fontWeight: questionType == QuestionType.newQuestion
                          ? FontWeight.bold
                          : null,
                      color: isDarkMode ? Colors.blue[200] : Colors.blue,
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    'Due: $amountDueQuestions',
                    style: TextStyle(
                      fontWeight: questionType == QuestionType.dueQuestion
                          ? FontWeight.bold
                          : null,
                      color: isDarkMode ? Colors.red[200] : Colors.red,
                    ),
                  ),
                ),
                Padding(
                  padding: const EdgeInsets.all(8.0),
                  child: Text(
                    'Non Due: $amountNonDueQuestions',
                    style: TextStyle(
                      fontWeight: questionType == QuestionType.nonDueQuestion
                          ? FontWeight.bold
                          : null,
                      color: isDarkMode ? Colors.green[200] : Colors.green,
                    ),
                  ),
                ),
              ],
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: shouldShowQuestion.value
                  ? [
                      Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: ElevatedButton(
                          onPressed: () {
                            shouldShowQuestion.value = false;
                          },
                          child: Text('Check'),
                        ),
                      )
                    ]
                  : [
                      Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: ElevatedButton(
                          onPressed: () {
                            try {
                              ref
                                  .read(quizQuestionsProvider.notifier)
                                  .gradeQuestion(
                                      remoteServer, tagsSet, id, true);
                            } on ServerException catch (e) {
                              if (context.mounted) {
                                ScaffoldMessenger.of(context).showSnackBar(
                                  SnackBar(
                                    backgroundColor: Colors.red,
                                    content: Text('$e'),
                                  ),
                                );
                              }
                            }
                          },
                          style: ButtonStyle(
                            backgroundColor:
                                WidgetStatePropertyAll(Colors.green),
                          ),
                          child: Text('Correct'),
                        ),
                      ),
                      Padding(
                        padding: const EdgeInsets.all(8.0),
                        child: ElevatedButton(
                          onPressed: () {
                            ref
                                .read(quizQuestionsProvider.notifier)
                                .gradeQuestion(
                                    remoteServer, tagsSet, id, false);
                          },
                          style: ButtonStyle(
                            backgroundColor: WidgetStatePropertyAll(Colors.red),
                          ),
                          child: Text('Incorrect'),
                        ),
                      )
                    ],
            ),
          ],
        ),
      ),
    );
  }
}
