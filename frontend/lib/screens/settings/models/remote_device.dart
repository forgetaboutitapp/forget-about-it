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
  const RemoteDevice({
    required this.title,
    required this.dateAdded,
    required this.lastUsed,
  });
  factory RemoteDevice.fromJSON(dynamic json) {
    final lastUsed = json['last-used'];
    return RemoteDevice(
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
