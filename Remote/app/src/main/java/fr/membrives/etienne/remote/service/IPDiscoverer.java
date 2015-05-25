package fr.membrives.etienne.remote.service;

import android.content.Context;
import android.net.wifi.WifiManager;
import android.util.Log;
import android.util.Pair;

import com.google.common.collect.Lists;
import com.google.protobuf.InvalidProtocolBufferException;

import org.zeromq.ZMQ;

import java.nio.ByteOrder;
import java.util.Formatter;
import java.util.List;
import java.util.Set;
import java.util.Timer;
import java.util.TimerTask;
import java.util.concurrent.Executor;

import fr.membrives.etienne.remote.RemoteProtos;

/**
 * IPDiscoverer discovers new peers on the WiFi IP subnet.
 */
public class IPDiscoverer implements IDiscoverer {
    private static final String ID_FORMAT = "tcp://%s:7001";
    private static final byte[] PING_MESSAGE =
            RemoteProtos.Request.newBuilder().setType(RemoteProtos.RequestType.PING).build()
                    .toByteArray();
    private final Context androidContext;
    private final ZMQ.Context zmqContext;
    private final Executor callerExecutor;
    private final Handler handler;
    private Timer discoveryTimer;
    private boolean started = false;

    private IPDiscoverer(Context androidContext, ZMQ.Context zmqContext, Handler handler,
                         Executor callerExecutor) {
        this.androidContext = androidContext;
        this.zmqContext = zmqContext;
        this.handler = handler;
        this.callerExecutor = callerExecutor;
    }

    private static String formatId(String... parameters) {
        return new Formatter().format(ID_FORMAT, parameters).out().toString();
    }

    @Override
    public void startDiscovery() {
        if (started) {
            return;
        }
        started = true;
        discoveryTimer = new Timer();
        discoveryTimer.schedule(new Discover(), 0, 1000 * 30);
    }

    @Override
    public void stopDiscovery() {
        started = false;
        discoveryTimer.cancel();
    }

    public static class Builder implements IPDiscoverer.Factory {
        private Context androidContext;
        private ZMQ.Context zmqContext;
        private Handler handler;
        private Executor callerExecutor;

        @Override
        public IDiscoverer createDiscoverer() {
            return new IPDiscoverer(this.androidContext, this.zmqContext, this.handler,
                    this.callerExecutor);
        }

        @Override
        public Factory setAndroidContext(Context androidContext) {
            this.androidContext = androidContext;
            return this;
        }

        @Override
        public Factory setZMQContext(ZMQ.Context zmqContext) {
            this.zmqContext = zmqContext;
            return this;
        }

        @Override
        public Factory setHandler(Handler handler) {
            this.handler = handler;
            return this;
        }

        @Override
        public Factory setCallerExecutor(Executor callerExecutor) {
            this.callerExecutor = callerExecutor;
            return this;
        }
    }

    private class Discover extends TimerTask {
        @Override
        public void run() {
            Set<String> ips = getAllIpAddressesInSubnet();
            List<Pair<String, ZMQ.Socket>> socketList = Lists.newArrayList();
            ZMQ.Poller poller = new ZMQ.Poller(ips.size());
            for (String ip : ips) {
                ZMQ.Socket socket = zmqContext.socket(ZMQ.REQ);
                socket.setReceiveTimeOut(10000);  // 10s timeout
                socket.setLinger(10000);
                socket.connect("tcp://" + ip);
                socket.send(PING_MESSAGE);
                poller.register(socket);
                socketList.add(Pair.create(ip, socket));
            }

            poller.poll(10000);

            for (int i = 0; i < socketList.size(); i++) {
                final Pair<String, ZMQ.Socket> pair = socketList.get(i);
                if (poller.pollin(i)) {
                    byte[] data = pair.second.recv(ZMQ.DONTWAIT);
                    if (data != null) {
                        final RemoteProtos.Response response;
                        try {
                            response = RemoteProtos.Response.parseFrom(data);
                            callerExecutor.execute(new Runnable() {
                                @Override
                                public void run() {
                                    IPDiscoverer.this.handler
                                            .newPeer(response.getHost(), formatId(pair.first));
                                }
                            });
                        } catch (InvalidProtocolBufferException e) {
                            Log.e("IPDiscoverer", "Unable to parse response from ip " + pair.first,
                                    e);
                        }
                    }
                } else {
                    callerExecutor.execute(new Runnable() {
                        @Override
                        public void run() {
                            IPDiscoverer.this.handler.disconnectedPeer(formatId(pair.first));
                        }
                    });
                }
                pair.second.close();
            }
        }

        private Set<String> getAllIpAddressesInSubnet() {
            WifiManager wifiManager =
                    (WifiManager) androidContext.getSystemService(Context.WIFI_SERVICE);
            int ipAddress = wifiManager.getConnectionInfo().getIpAddress();

            // Convert little-endian to big-endianif needed
            if (ByteOrder.nativeOrder().equals(ByteOrder.LITTLE_ENDIAN)) {
                ipAddress = Integer.reverseBytes(ipAddress);
            }

            int netmask = wifiManager.getDhcpInfo().netmask;

            // Convert little-endian to big-endianif needed
            if (ByteOrder.nativeOrder().equals(ByteOrder.LITTLE_ENDIAN)) {
                netmask = Integer.reverseBytes(netmask);
            }

            return IPUtils.getIpsInSubnet(ipAddress, netmask, 0);
        }
    }
}
