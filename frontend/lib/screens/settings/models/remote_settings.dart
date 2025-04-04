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

  static RemoteSettings fromJSON(dynamic json) {
    final List<RemoteDevice>? remoteDevices = json['settings']
            ?['remote-devices']
        ?.map<RemoteDevice>((e) => RemoteDevice.fromJSON(e))
        ?.toList();
    final List<RemoteAlgorithm>? remoteAlgorithms = json['settings']
            ?['remote-algorithms']
        ?.map<RemoteAlgorithm>((e) => RemoteAlgorithm.fromJSON(e))
        ?.toList();
    final iRemoteAlgorithm = remoteAlgorithms?.toIList();
    final int? def = iRemoteAlgorithm?.firstOrNull?.algorithmID;
    final sortedRemoteAlgorithm = iRemoteAlgorithm
        ?.sort((a, b) => a.algorithmID.compareTo(b.algorithmID));

    return RemoteSettings(
      remoteDevices: remoteDevices?.toIList(),
      remoteAlgorithms: sortedRemoteAlgorithm?.sort(),
      defaultAlgorithm: def,
    );
  }
}
