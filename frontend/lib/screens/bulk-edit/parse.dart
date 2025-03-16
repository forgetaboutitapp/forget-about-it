import 'dart:collection';

import 'package:app/screens/bulk-edit/model.dart';
import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:flutter/material.dart';

String unparse(IList<Flashcard> flashcards) => flashcards
    .map((e) =>
        '${_decode(e.id)}|${_decode(e.question)}|${e.answer}|${e.tags.map((e) => _decode(e)).join(' ')}')
    .join('\n');

String _decode<A>(A a) => a
    .toString()
    .replaceAll('\n', '\\n')
    .replaceAll('|', '\\|')
    .replaceAll('\\', '\\\\')
    .replaceAll('%', '\\%');

IList<Flashcard> parse(String data) {
  final splitData = data.split('\n');
  final flashcardList = List.generate(splitData.length, (index) {
    final line = removeCommentsAndSplit(index, splitData[index]);
    if (line.length == 1 && line[0] == '') {
      return null;
    }
    if (line.length != 4) {
      throw InvalidNumberOfFields(
          lineNumber: index,
          stringWithError: splitData[index],
          numberOfFields: line.length);
    }
    int? id;
    try {
      id = line[0] == '' ? null : int.parse(line[0]);
    } on FormatException catch (_) {
      throw QuestionIDNotIntException(
          id: line[0], question: line[1], lineNumber: index);
    }
    return Flashcard(
        id: id,
        question: line[1],
        answer: line[2],
        tags: line[3].split(' ').where((e) => e.trim() != '').toIList());
  }).where((f) => f != null).map((e) => e as Flashcard).toIList();
  // ensure that ids don't conflic and ensure all questions have tags
  HashSet cardIDs = HashSet();
  HashSet cardQuestions = HashSet();

  for (final flashcard in flashcardList) {
    if (flashcard.tags.isEmpty ||
        (flashcard.tags.length == 1 && flashcard.tags[0] == '')) {
      throw NoTagException(id: flashcard.id, question: flashcard.question);
    }
    final id = flashcard.id;
    final question = flashcard.question;
    if (id != null) {
      if (cardIDs.contains(id)) {
        throw IDsConflict(
          conflictingID: id,
        );
      }
      cardIDs.add(id);
    }
    if (cardQuestions.contains(question)) {
      throw QuestionsConflict(conflictingQuestion: question);
    }
    cardQuestions.add(question);
  }

  return flashcardList;
}

IList<String> removeCommentsAndSplit(int lineNumber, String splitData) {
  IList<String> returnVal = const IList.empty();
  final curString = StringBuffer();
  final stringIterator = splitData.characters;
  bool shouldSkip = false;
  loop:
  for (final c in stringIterator) {
    if (shouldSkip) {
      if (c == 'n') {
        curString.write('\n');
      } else {
        curString.write(c);
      }
      shouldSkip = false;
    } else {
      switch (c) {
        case '\\':
          {
            shouldSkip = true;
          }
        case '|':
          {
            returnVal = returnVal.add(curString.toString().trim());
            curString.clear();
          }
        case '%':
          {
            break loop;
          }
        default:
          {
            curString.write(c);
          }
      }
    }
  }
  if (shouldSkip) {
    throw InvalidEscapeException(
        lineNumber: lineNumber, stringWithError: splitData);
  }
  returnVal = returnVal.add(curString.toString().trim());
  return returnVal;
}

class InvalidEscapeException implements Exception {
  final int lineNumber;
  final String stringWithError;

  InvalidEscapeException(
      {required this.stringWithError, required this.lineNumber});
  @override
  String toString() => 'Invalid escape on line $lineNumber: $stringWithError';
}

class InvalidNumberOfFields implements Exception {
  final int lineNumber;
  final int numberOfFields;
  final String stringWithError;

  InvalidNumberOfFields(
      {required this.stringWithError,
      required this.lineNumber,
      required this.numberOfFields});

  @override
  String toString() =>
      'Invalid number ($numberOfFields) of fields on line $lineNumber: $stringWithError';
}

class IDsConflict implements Exception {
  final int conflictingID;

  IDsConflict({required this.conflictingID});
  @override
  String toString() => 'Two questions have the same ID ($conflictingID)';
}

class QuestionsConflict implements Exception {
  final String conflictingQuestion;

  QuestionsConflict({required this.conflictingQuestion});
  @override
  String toString() =>
      'Two questions have the same Question ($conflictingQuestion)';
}

class NoTagException implements Exception {
  final int? id;
  final String question;

  NoTagException({required this.id, required this.question});
  @override
  String toString() => 'Question $id ($question) does not contain tags';
}

class QuestionIDNotIntException implements Exception {
  final String? id;
  final String question;
  final int lineNumber;

  QuestionIDNotIntException({
    required this.id,
    required this.question,
    required this.lineNumber,
  });
  @override
  String toString() => 'Question id $id is not an integer';
}
