package fr.membrives.etienne.commandforwarder.service;

import android.util.Pair;

import com.google.common.util.concurrent.AsyncFunction;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;

import org.zeromq.ZMQ;

import java.util.concurrent.Callable;
import java.util.concurrent.Executors;

import fr.membrives.etienne.remote.RemoteProtos;

/**
 * Main forwarding service.
 */
public class ForwarderService {
    private ListeningExecutorService executor = MoreExecutors.listeningDecorator(Executors.newFixedThreadPool(3));

    private ZMQ.Context mContext;
    private ZMQ.Socket mSendingSocket;
    private ZMQ.Socket mControlSocket;

    public int ForwarderService() {
        mContext = ZMQ.context(1);
        mControlSocket = mContext.socket(ZMQ.REQ);
        mControlSocket.bind("tcp://etienne.membrives.fr:6003");

        mSendingSocket = mContext.socket(ZMQ.PUB);
        mSendingSocket.bind("tcp://etienne.membrives.fr:6002");
    }

    public void sendCommand(String command) {
        executor.execute(getCommandSender(command));
    }

    /**
     * @return the current status of (server, frontend)
     */
    public ListenableFuture<Pair<Boolean, Boolean>> getStatus() {
        if (mControlSocket.)
        ListenableFuture<RemoteProtos.Response> future = executor.submit(getStatusCallable());
        Futures.transform(future, new AsyncFunction<RemoteProtos.Response, Pair<Boolean, Boolean>>() {
            @Override
            public ListenableFuture<Pair<Boolean, Boolean>> apply(RemoteProtos.Response input) throws Exception {
                return Pair.of(true, input.getStatus().equals(RemoteProtos.Response.FrontendStatus.CONNECTED));
            }
        }, executor);
    }

    private Runnable getCommandSender(final String command) {
        return new Runnable() {
            @Override
            public void run() {
                RemoteProtos.Command commandProto = RemoteProtos.Command.newBuilder()
                        .setSource("android").setType(RemoteProtos.Command.CommandType.COMMAND)
                        .setCommand(command).build();
                mSendingSocket.send(commandProto.toByteArray());
            }
        };
    }

    private Callable<RemoteProtos.Response> getStatusCallable() {
        return new Callable<RemoteProtos.Response>() {
            @Override
            public RemoteProtos.Response call() throws Exception {
                RemoteProtos.Command commandProto = RemoteProtos.Command.newBuilder()
                        .setSource("android").setType(RemoteProtos.Command.CommandType.STATUS).build();
                mControlSocket.send(commandProto.toByteArray());
            }
        }
    }
}
