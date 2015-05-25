package fr.membrives.etienne.remote;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.view.KeyEvent;
import android.view.WindowManager;
import android.widget.ImageView;

import fr.membrives.etienne.remote.service.ForwarderService;


public class RemoteController extends Activity {
    private static final String TAG = "RemoteController";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_remote_controller);
        getWindow().addFlags(WindowManager.LayoutParams.FLAG_KEEP_SCREEN_ON);
        ImageView serverCheck = (ImageView) findViewById(R.id.server_img);

        // Start the forwarder service
        Intent i = new Intent(this, ForwarderService.class);
        this.startService(i);
    }

    @Override
    public boolean onKeyUp(int keyCode, KeyEvent event) {
        return super.onKeyUp(keyCode, event);
    }
}
