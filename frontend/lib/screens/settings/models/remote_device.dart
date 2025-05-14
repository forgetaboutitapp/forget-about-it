import 'package:freezed_annotation/freezed_annotation.dart';

part 'remote_device.freezed.dart';

@freezed
class RemoteDevice with _$RemoteDevice {
  @override
  final DateTime dateAdded;
  @override
  final DateTime? lastUsed;
  @override
  final String title;
  @override
  final String loginId;
  const RemoteDevice({
    required this.title,
    required this.dateAdded,
    required this.lastUsed,
    required this.loginId,
  });
}
