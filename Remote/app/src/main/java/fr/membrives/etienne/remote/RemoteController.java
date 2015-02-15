package fr.membrives.etienne.remote;

import android.app.Activity;
import android.os.Bundle;
import android.util.Log;
import android.view.KeyEvent;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.view.WindowManager;
import android.widget.Button;
import android.widget.ImageView;

import java.util.concurrent.ExecutionException;

import fr.membrives.etienne.remote.service.ForwarderService;


public class RemoteController extends Activity {
    private static final String TAG = "f.m.e.remote.RemoteController";
    private ForwarderService service;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_remote_controller);
        getWindow().addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);
        ImageView serverCheck = (ImageView) findViewById(R.id.server_img);
        service = new ForwarderService();
        try {
            if (service.connect(this).get()) {
                serverCheck.setImageResource(R.drawable.check);
            } else {
                serverCheck.setImageResource(R.drawable.fail);
            }
        } catch (InterruptedException e) {
            Log.e(TAG, "Interrupted", e);
        } catch (ExecutionException e) {
            Log.e(TAG, "Interrupted", e);
        }

        final Button prevBtn = (Button) findViewById(R.id.prev_btn);
        prevBtn.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                RemoteProtos.Command.Builder commandBuilder = RemoteProtos.Command.newBuilder().setType(RemoteProtos.Command.CommandType.COMMAND);
                commandBuilder.setCommand("p");
                service.sendWebcontrolMessage(commandBuilder.build());
            }
        });
        final Button nextBtn = (Button) findViewById(R.id.next_btn);
        nextBtn.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
                RemoteProtos.Command.Builder commandBuilder = RemoteProtos.Command.newBuilder().setType(RemoteProtos.Command.CommandType.COMMAND);
                commandBuilder.setCommand("n");
                service.sendWebcontrolMessage(commandBuilder.build());
            }
        });
    }


    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        // Inflate the menu; this adds items to the action bar if it is present.
        getMenuInflater().inflate(R.menu.menu_remote_controller, menu);
        return true;
    }

    @Override
    public boolean onKeyUp(int keyCode, KeyEvent event) {
        RemoteProtos.Command.Builder commandBuilder = RemoteProtos.Command.newBuilder().setType(RemoteProtos.Command.CommandType.COMMAND);
        commandBuilder.setCommand(String.valueOf((char)event.getUnicodeChar()));
        service.sendWebcontrolMessage(commandBuilder.build());
        return super.onKeyUp(keyCode, event);
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {
        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.
        int id = item.getItemId();

        //noinspection SimplifiableIfStatement
        if (id == R.id.action_settings) {
            return true;
        }

        return super.onOptionsItemSelected(item);
    }
}
