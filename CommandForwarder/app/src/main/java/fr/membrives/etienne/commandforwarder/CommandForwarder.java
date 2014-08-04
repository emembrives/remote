package fr.membrives.etienne.commandforwarder;

import android.app.Activity;
import android.os.Bundle;
import android.view.KeyEvent;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.Toast;

import fr.membrives.etienne.commandforwarder.service.ForwarderService;
import fr.membrives.etienne.remote.RemoteProtos;
import nanomsg.exceptions.IOException;


public class CommandForwarder extends Activity {
    private ForwarderService service;


    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_command_forwarder);
        service = new ForwarderService();
    }

    @Override
    protected void onDestroy() {
        service.stop();
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
        commandBuilder.setCommand(event.getKeyCode());
        try {
            service.submit(commandBuilder.build());
        } catch (IOException e) {
            Toast.makeText(this, "Error sending command", Toast.LENGTH_LONG).show();
        }
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

    public void updateDisplay() {

    }
}
