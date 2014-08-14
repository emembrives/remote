package fr.membrives.etienne.commandforwarder;

import android.app.Activity;
import android.os.Bundle;
import android.util.Log;
import android.view.KeyEvent;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.ImageView;

import java.io.IOException;
import java.util.concurrent.ExecutionException;

import fr.membrives.etienne.commandforwarder.service.ForwarderService;
import fr.membrives.etienne.remote.RemoteProtos;


public class CommandForwarder extends Activity {
    private static final String TAG = "fr.m.e.cf.CommandForwarder";

    private ForwarderService service;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_command_forwarder);
        ImageView serverCheck = (ImageView) findViewById(R.id.server_img);
        service = new ForwarderService();
        try {
            if (service.connect().get()) {
                serverCheck.setImageResource(R.drawable.check);
            } else {
                serverCheck.setImageResource(R.drawable.fail);
            }
        } catch (InterruptedException e) {
            Log.e(TAG, "Interrupted", e);
        } catch (ExecutionException e) {
            Log.e(TAG, "Interrupted", e);
        }
    }

    @Override
    protected void onDestroy() {
        try {
            service.stopForwarderService();
        } catch (IOException e) {
            Log.e(TAG, "Failed to close connection", e);
        }
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.command_forwarder, menu);
        return true;
    }

    @Override
    public boolean onKeyUp(int keyCode, KeyEvent event) {
        RemoteProtos.Command.Builder commandBuilder = RemoteProtos.Command.newBuilder().setType(RemoteProtos.Command.CommandType.COMMAND);
        commandBuilder.setCommand(event.getCharacters());
        service.sendMessage(commandBuilder.getCommandBytes());
        return super.onKeyUp(keyCode, event);
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();
        if (id == R.id.action_settings) {
            return true;
        }
        return super.onOptionsItemSelected(item);
    }
}
