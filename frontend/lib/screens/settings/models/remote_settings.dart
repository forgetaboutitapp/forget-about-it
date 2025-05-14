import 'package:fast_immutable_collections/fast_immutable_collections.dart';
import 'package:freezed_annotation/freezed_annotation.dart';

import 'remote_algorithm.dart';
import 'remote_device.dart';

part 'remote_settings.freezed.dart';

@freezed
class RemoteSettings with _$RemoteSettings {
  const RemoteSettings({
    required this.remoteDevices,
    required this.remoteAlgorithms,
    required this.defaultAlgorithm,
  });
  @override
  final IList<RemoteDevice>? remoteDevices;
  @override
  final IList<RemoteAlgorithm>? remoteAlgorithms;
  @override
  final int? defaultAlgorithm;
}
