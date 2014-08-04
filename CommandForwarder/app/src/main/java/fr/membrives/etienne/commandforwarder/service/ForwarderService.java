package fr.membrives.etienne.commandforwarder.service;

import com.google.common.util.concurrent.AsyncFunction;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListenableFutureTask;
import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;
import com.google.protobuf.InvalidProtocolBufferException;

import java.util.concurrent.Callable;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.LinkedBlockingQueue;

import fr.membrives.etienne.remote.RemoteProtos;
import nanomsg.Message;
import nanomsg.exceptions.IOException;
import nanomsg.pair.PairSocket;

/**
 * Main forwarding service.
 */
public class ForwarderService implements Runnable {
    public static final AsyncFunction<RemoteProtos.Response, Boolean> RESPONSE_TO_FRONTEND_CONNECTED = new AsyncFunction<RemoteProtos.Response, Boolean>() {
        @Override
        public ListenableFuture<Boolean> apply(RemoteProtos.Response input) throws Exception {
            return Futures.immediateCheckedFuture(input.getFrontendConnected());
        }
    };
    private boolean isStopped = true;

    private boolean isConnected = false;

    private LinkedBlockingQueue<RemoteProtos.Response> responses;
    private ListeningExecutorService executor = MoreExecutors.listeningDecorator(Executors.newFixedThreadPool(3));

    private PairSocket socket;

    public void stop() {
        isStopped = true;
        executor.shutdownNow();
    }

    public void start() {
        if (!isConnected) {
            executor.execute(this);
        }
    }

    @Override
    public void run() {
        try {
            socket.connect("tcp://polaris.membrives.fr:6001");
        } catch (IOException e) {
            return;
        }
        isConnected = true;
        while (!isStopped) {
            try {
                Message message = socket.recv();
                RemoteProtos.Response response = RemoteProtos.Response.parseFrom(message.toBytes());
                responses.offer(response);
            } catch (IOException e) {
                e.printStackTrace();
                return;
            } catch (InvalidProtocolBufferException e) {
                e.printStackTrace();
                return;
            }
        }
    }

    public void submit(final RemoteProtos.Command message) throws IOException {
        executor.submit(new Callable<Void>() {
            @Override
            public Void call() throws Exception {
                socket.sendBytes(message.toByteArray());
                return null;
            }
        });
    }

    public ListenableFuture<Boolean> isFrontendConnected() throws IOException {
        RemoteProtos.Command.Builder builder = RemoteProtos.Command.newBuilder().setType(RemoteProtos.Command.CommandType.STATUS);
        submit(builder.build());
        ListenableFuture<RemoteProtos.Response> msgFuture = executor.submit(new Callable<RemoteProtos.Response>() {
            @Override
            public RemoteProtos.Response call() throws Exception {
                return responses.take();
            }
        });
        return Futures.transform(msgFuture, RESPONSE_TO_FRONTEND_CONNECTED);
    }

    public boolean isConnected() {
        return isConnected;
    }
}
