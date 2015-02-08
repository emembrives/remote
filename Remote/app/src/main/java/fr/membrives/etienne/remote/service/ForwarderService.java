package fr.membrives.etienne.remote.service;

import android.content.Context;
import android.util.Log;

import com.google.common.util.concurrent.FutureCallback;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;

import java.io.IOException;
import java.util.concurrent.Callable;
import java.util.concurrent.Executors;

import fr.membrives.etienne.remote.R;
import fr.membrives.etienne.remote.RemoteProtos;
import nanomsg.reqrep.ReqSocket;

/**
 * Main forwarding service.
 */
public class ForwarderService {
    private static final String TAG = "f.m.e.remote.service.ForwarderService";
    private static final String EXCHANGE_NAME = "remote";

    private ListeningExecutorService executor = MoreExecutors.listeningDecorator(Executors.newFixedThreadPool(1));
    private boolean connectedToServer = false;
    private ReqSocket socket;

    public ForwarderService() {
    }

    public ListenableFuture<Boolean> connect(final Context context) {
        ListenableFuture<Boolean> serverConnect = executor.submit(new Callable<Boolean>() {
            @Override
            public Boolean call() {
                socket = new ReqSocket();
                try {
                    socket.connect("tcp://terra.membrives.fr:7001");
                } catch (nanomsg.exceptions.IOException e) {
                    Log.e(TAG, "Unable to connect", e);
                    return false;
                }
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
                try {
                    socket.sendBytes(message.toByteArray());
                    return true;
                } catch (IOException e) {
                    Log.e(TAG, "Unable to send message", e);
                    return false;
                }
            }
        });
        return messageSent;
    }

    public void stopForwarderService() throws IOException {
        socket.close();
        connectedToServer = false;
    }
}
