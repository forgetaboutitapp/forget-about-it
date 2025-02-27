import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:flutter/foundation.dart';
part 'keys.freezed.dart';

@freezed
class TwelveWordKeys with _$TwelveWordKeys {
  const factory TwelveWordKeys({required int row, required int col}) =
      _TwelveWordKeys;
}

const twelveWordLoginButton = 10001;
const loginURLKey = 10002;
