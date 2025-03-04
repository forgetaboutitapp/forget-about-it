import 'package:fast_immutable_collections/fast_immutable_collections.dart';

sealed class SubmitType {}

class TwelveWords extends SubmitType {
  final IList<String> twelveWords;

  TwelveWords({required this.twelveWords});
}

class Token extends SubmitType {
  final String token;

  Token({required this.token});
}
