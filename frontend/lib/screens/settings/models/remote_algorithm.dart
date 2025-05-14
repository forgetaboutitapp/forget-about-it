import 'package:freezed_annotation/freezed_annotation.dart';

part 'remote_algorithm.freezed.dart';

@freezed
class RemoteAlgorithm with _$RemoteAlgorithm {
  const RemoteAlgorithm({
    required this.algorithmID,
    required this.authorName,
    required this.license,
    required this.remoteURL,
    required this.downloadURL,
    required this.version,
    required this.algorithmName,
    required this.timeAdded,
  });
  @override
  final int algorithmID;
  @override
  final String authorName;
  @override
  final String license;
  @override
  final String remoteURL;
  @override
  final String downloadURL;
  @override
  final String timeAdded;
  @override
  final int version;
  @override
  final String algorithmName;
}
