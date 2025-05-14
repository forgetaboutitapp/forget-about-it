import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'model.freezed.dart';

@freezed
class Flashcard with _$Flashcard {
  const Flashcard({
    required this.id,
    required this.question,
    required this.answer,
    required this.explanation,
    required this.memoHint,
    required this.tags,
  });
  @override
  final int? id;
  @override
  final String question;
  @override
  final String answer;
  @override
  final String explanation;
  @override
  final String memoHint;
  @override
  final IList<String> tags;
}
