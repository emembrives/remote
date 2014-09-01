package fr.membrives.etienne.commandforwarder.service;

import android.content.Context;
import android.util.Log;

import com.google.common.util.concurrent.FutureCallback;
import com.google.common.util.concurrent.Futures;
import com.google.common.util.concurrent.ListenableFuture;
import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;
import com.google.protobuf.ByteString;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;

import java.io.IOException;
import java.util.concurrent.Callable;
import java.util.concurrent.Executors;

import fr.membrives.etienne.commandforwarder.R;

/**
 * Main forwarding service.
 */
public class ForwarderService {
    private static final String TAG = "fme.c.s.ForwarderService";
    private static final String EXCHANGE_NAME = "remote";

    private Channel channel;
    private Connection connection;
    private ListeningExecutorService executor = MoreExecutors.listeningDecorator(Executors.newFixedThreadPool(1));
    private boolean connectedToServer = false;

    public ForwarderService() {
    }

    public ListenableFuture<Boolean> connect(final Context context) {
        ListenableFuture<Boolean> serverConnect = executor.submit(new Callable<Boolean>() {
            @Override
            public Boolean call() {
                ConnectionFactory factory = new ConnectionFactory();
                factory.setHost(context.getString(R.string.amqp_server));
                factory.setUsername(context.getString(R.string.amqp_user));
                factory.setPassword(context.getString(R.string.amqp_password));
                try {
                    connection = factory.newConnection();
                    channel = connection.createChannel();
                    channel.exchangeDeclare(EXCHANGE_NAME, "topic", true, true, false, null);
                } catch (IOException e) {
                    Log.e(TAG, "Not connected to server", e);
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

    public ListenableFuture<Boolean> sendWebcontrolMessage(final ByteString message) {
        ListenableFuture<Boolean> messageSent = executor.submit(new Callable<Boolean>() {
            @Override
            public Boolean call() {
                try {
                    channel.basicPublish(EXCHANGE_NAME, "webcontrol", null, message.toByteArray());
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
        channel.close();
        connection.close();
    }
}
