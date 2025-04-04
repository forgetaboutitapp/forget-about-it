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
  factory RemoteDevice.fromJSON(dynamic json) {
    final lastUsed = json['last-used'];
    return RemoteDevice(
      loginId: json['login-id'],
      title: json['title'],
      dateAdded: DateTime.fromMillisecondsSinceEpoch(
        json['date-added'] * 1000,
      ),
      lastUsed: (lastUsed == null)
          ? null
          : DateTime.fromMillisecondsSinceEpoch(
              lastUsed * 1000,
            ),
    );
  }
}
