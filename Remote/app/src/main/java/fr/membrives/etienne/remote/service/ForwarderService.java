package fr.membrives.etienne.remote.service;

import android.content.Context;

import com.google.common.util.concurrent.FutureCallback;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;

import org.zeromq.ZMQ;

import java.io.IOException;
import java.util.concurrent.Callable;
import java.util.concurrent.Executors;

import fr.membrives.etienne.remote.RemoteProtos;

/**
 * Main forwarding service.
 */
public class ForwarderService {
    private static final String TAG = "s.ForwarderService";
    private static final String EXCHANGE_NAME = "remote";

    private ListeningExecutorService executor = MoreExecutors.listeningDecorator(Executors.newFixedThreadPool(1));
    private boolean connectedToServer = false;

    private ZMQ.Context zmqContext;
    private ZMQ.Socket socket;

    public ForwarderService() {
    }

    public ListenableFuture<Boolean> connect(final Context context) {
        ListenableFuture<Boolean> serverConnect = executor.submit(new Callable<Boolean>() {
            @Override
            public Boolean call() {
                zmqContext = ZMQ.context(1);
                socket = zmqContext.socket(ZMQ.REQ);
                socket.connect("tcp://10.0.2.2:7002");
                return true;
            }
        });
        Futures.addCallback(serverConnect, new FutureCallback<Boolean>() {
            @Override
            public void onSuccess(Boolean result) {
                connectedToServer = result;
            }

            @Override
            public void onFailure(Throwable t) {
                connectedToServer = false;
            }
        });
        return serverConnect;
    }

    public ListenableFuture<Boolean> sendWebcontrolMessage(final RemoteProtos.Command message) {
        ListenableFuture<Boolean> messageSent = executor.submit(new Callable<Boolean>() {
            @Override
            public Boolean call() {
                return socket.send(message.toByteArray(), 0);
            }
        });
        return messageSent;
    }

    public void stopForwarderService() throws IOException {
        socket.close();
        zmqContext.term();
        connectedToServer = false;
    }
}
