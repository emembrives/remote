package fr.membrives.etienne.remote.service;

import android.app.Service;
import android.content.Intent;
import android.os.Binder;
import android.os.IBinder;
import android.util.Log;

import com.google.common.collect.Lists;
import com.google.protobuf.InvalidProtocolBufferException;

import org.zeromq.ZMQ;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;

import fr.membrives.etienne.remote.RemoteProtos;

/**
 * Main forwarding service.
 */
public class ForwarderService extends Service implements IDiscoverer.Handler {
    private static final String TAG = "s.ForwarderService";
    private final IBinder mBinder = new ForwarderServiceBinder();
    private ZMQ.Context zmqContext;
    private List<IDiscoverer> discovererList = Lists.newArrayList();
    private Executor executor = Executors.newSingleThreadExecutor();
    private Map<String, ZMQPeer> peers = new HashMap<String, ZMQPeer>();
    public ForwarderService() {
    }

    @Override
    public int onStartCommand(Intent intent, int flags, int startId) {
        zmqContext = ZMQ.context(1);
        discovererList
                .add(new IPDiscoverer.Builder().setAndroidContext(this).setZMQContext(zmqContext)
                        .setHandler(this).setCallerExecutor(executor).createDiscoverer());
        discoverPeers();
        return Service.START_NOT_STICKY;
    }

    @Override
    public void onDestroy() {
        for (IDiscoverer discoverer : discovererList) {
            discoverer.stopDiscovery();
        }
        zmqContext.term();
    }

    @Override
    public IBinder onBind(Intent intent) {
        return mBinder;
    }

    public void discoverPeers() {
        for (IDiscoverer discoverer : discovererList) {
            discoverer.startDiscovery();
        }
    }

    @Override
    public void newPeer(String host, String id) {
        if (peers.containsKey(id)) {
            return;
        }
        ZMQ.Socket socket = zmqContext.socket(ZMQ.REQ);
        socket.connect(id);
        ZMQPeer peer = new ZMQPeer(socket, host, id);
        peers.put(id, peer);
    }

    @Override
    public void disconnectedPeer(String id) {
        if (peers.containsKey(id)) {
            ZMQPeer peer = peers.get(id);
            peer.socket.close();
            peers.remove(id);
        }
    }

    public static class ZMQPeer {
        private final ZMQ.Socket socket;
        private final String host;
        private final String id;

        private List<RemoteProtos.ServiceDefinition> services = Lists.newArrayList();

        private ZMQPeer(ZMQ.Socket socket, String host, String id) {
            this.socket = socket;
            this.host = host;
            this.id = id;
        }

        public void getServices() {
            RemoteProtos.Request.Builder builder = RemoteProtos.Request.newBuilder();
            builder.setHost(host);
            builder.setType(RemoteProtos.RequestType.SERVICES);
            this.socket.send(builder.build().toByteArray());
            byte[] responseBytes = this.socket.recv();
            RemoteProtos.Response response;
            try {
                response = RemoteProtos.Response.parseFrom(responseBytes);
            } catch (InvalidProtocolBufferException e) {
                Log.e("ForwarderService", "Unable to parse services response for host " + host, e);
                return;
            }
            services = response.getServicesList();
        }
    }

    public class ForwarderServiceBinder extends Binder {
        ForwarderService getService() {
            return ForwarderService.this;
        }
    }
}
