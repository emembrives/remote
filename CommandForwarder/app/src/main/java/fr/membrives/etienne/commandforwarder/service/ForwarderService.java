package fr.membrives.etienne.commandforwarder.service;

import com.google.common.util.concurrent.ListeningExecutorService;
import com.google.common.util.concurrent.MoreExecutors;
import com.rabbitmq.client.Channel;
import com.rabbitmq.client.Connection;
import com.rabbitmq.client.ConnectionFactory;

import java.io.IOException;
import java.util.concurrent.Executors;

/**
 * Main forwarding service.
 */
public class ForwarderService {
    private final static String QUEUE_NAME = "hello";

    private Channel channel;
    private Connection connection;

    public ForwarderService() throws IOException {
        ConnectionFactory factory = new ConnectionFactory();
        factory.setHost("localhost");
        factory.setUsername("guest");
        factory.setPassword("guest");
        connection = factory.newConnection();
        channel = connection.createChannel();
        channel.queueDeclare(QUEUE_NAME, false, false, false, null);
    }

    public void sendMessage(String message) throws IOException {
        channel.basicPublish("", QUEUE_NAME, null, message.getBytes());
    }

    public void stopForwarderService() throws IOException {
        channel.close();
        connection.close();
    }
}
