import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'model.freezed.dart';

@freezed
class Flashcard with _$Flashcard {
  const Flashcard({
    required this.id,
    required this.question,
    required this.answer,
    required this.tags,
  });
  @override
  final int? id;
  @override
  final String question;
  @override
  final String answer;
  @override
  final IList<String> tags;

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'question': question,
      'answer': answer,
      'tags': tags.toList()
    };
  }

  static Flashcard fromJson(Map<String, dynamic> data) {
    int id = data['id'];
    String question = data['question'];
    String answer = data['answer'];
    final tags = List<String>.from(data['tags']);
    return Flashcard(
        id: id, question: question, answer: answer, tags: tags.toIList());
  }
}
