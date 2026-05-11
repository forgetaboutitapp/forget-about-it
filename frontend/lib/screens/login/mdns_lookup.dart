import 'package:multicast_dns/multicast_dns.dart';
import '../../interop/notification_service_stub.dart';

final Map<String, String> resolvedMdnsHosts = {};

Future<String?> discoverServerViaMdns() async {
  final MDnsClient client = MDnsClient();
  try {
    await client.start();
    String? foundUrl;

    await for (final PtrResourceRecord ptr in client.lookup<PtrResourceRecord>(
        ResourceRecordQuery.serverPointer('_forgetaboutit._tcp.local'))) {
      await for (final SrvResourceRecord srv
          in client.lookup<SrvResourceRecord>(
              ResourceRecordQuery.service(ptr.domainName))) {
        await for (final IPAddressResourceRecord ip
            in client.lookup<IPAddressResourceRecord>(
                ResourceRecordQuery.addressIPv4(srv.target))) {
          resolvedMdnsHosts[srv.target] = ip.address.address;
          foundUrl = 'http://${srv.target}:${srv.port}';
          break;
        }
        if (foundUrl != null) break;
      }
      if (foundUrl != null) break;
    }

    return foundUrl;
  } catch (e) {
    return null;
  } finally {
    client.stop();
  }
}

Future<void> resolveAndCacheMdnsHost(String mdnsHost) async {
  final MDnsClient client = MDnsClient();
  try {
    await client.start();
    await for (final IPAddressResourceRecord ip
        in client.lookup<IPAddressResourceRecord>(
            ResourceRecordQuery.addressIPv4(mdnsHost))) {
      resolvedMdnsHosts[mdnsHost] = ip.address.address;
      break;
    }
  } catch (e) {
    await NotificationService().showNotification(
      title: 'mDNS Lookup Error',
      body: 'Failed to resolve server IP: $e',
    );
  } finally {
    client.stop();
  }
}
