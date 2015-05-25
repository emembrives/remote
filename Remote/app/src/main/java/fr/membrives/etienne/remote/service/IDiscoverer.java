package fr.membrives.etienne.remote.service;

import android.content.Context;

import org.zeromq.ZMQ;

import java.util.concurrent.Executor;

public interface IDiscoverer {
    public void startDiscovery();

    public void stopDiscovery();

    public interface Factory {
        public Factory setAndroidContext(Context androidContext);
        public Factory setZMQContext(ZMQ.Context zmqContext);
        public Factory setHandler(Handler handler);
        public Factory setCallerExecutor(Executor callerExecutor);

        public IDiscoverer createDiscoverer();
    }

    public interface Handler {
        public void newPeer(String host, String id);

        public void disconnectedPeer(String id);
    }
}
