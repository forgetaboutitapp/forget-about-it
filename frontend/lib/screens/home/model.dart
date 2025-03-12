import 'package:freezed_annotation/freezed_annotation.dart';

part 'model.freezed.dart';

@freezed
class Tag with _$Tag {
  Tag({
    required this.tag,
    required this.totalQuestions,
  });

  @override
  String tag;

  @override
  int totalQuestions;

  static Tag fromJson(Map<String, dynamic> d) {
    return Tag(tag: d['tag'], totalQuestions: d['num-questions']);
  }
}
