import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

part 'submit_type.freezed.dart';

@freezed
sealed class SubmitType with _$SubmitType {
  const SubmitType._();
  factory SubmitType.twelveWords(IList<String> twelveWords) = TwelveWords;
  factory SubmitType.token(String token) = Token;
}
