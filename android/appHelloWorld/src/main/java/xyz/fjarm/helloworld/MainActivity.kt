package xyz.fjarm.helloworld

import android.os.Bundle
import android.util.Log
import android.widget.Toast
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.Dp
import androidx.lifecycle.lifecycleScope
import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import io.grpc.ManagedChannel
import io.grpc.StatusException
import io.grpc.android.AndroidChannelBuilder
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.asExecutor
import kotlinx.coroutines.launch
import xyz.fjarm.helloworld.ui.theme.HelloWorldTheme
import xyz.fjarm.libhelloworld.HelloWorldRepository
import xyz.fjarm.libhelloworld.HelloWorldGrpcClientImpl
import java.util.concurrent.TimeUnit

class MainActivity : ComponentActivity() {

    private lateinit var client: HelloWorldRepository

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            HelloWorldTheme {
                // A surface container using the 'background' color from the theme
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    Greeting(this::fetchGreeting)
                }
            }
        }
        client = HelloWorldGrpcClientImpl(initOkHttpChannel())
    }

    private fun initOkHttpChannel(): ManagedChannel {
        val channel = AndroidChannelBuilder
            // TODO(2024-07-11): Dagger inject an address instead of hardcoding localhost
            .forAddress("10.0.2.2", 8000)
            .context(applicationContext)
            .idleTimeout(30, TimeUnit.SECONDS)
            .keepAliveTime(2, TimeUnit.MINUTES)
            .executor(Dispatchers.IO.asExecutor())
            // TODO(2024-07-11): Remove this after securing the client with TLS certs
            .usePlaintext()
        return channel.build()
    }

    private fun fetchGreeting() {
        lifecycleScope.launch {
            try {
                val res = client
                    .getHelloWorld(GetHelloWorldRequest.newBuilder().build())
                Log.println(Log.INFO, "MainActivity", "jmuhia, fetchGreeting, res: $res")
                Toast.makeText(
                    this@MainActivity,
                    "Server returned response with output: ${res.output}",
                    Toast.LENGTH_LONG
                ).show()
            } catch (e: StatusException) {
                Log.println(Log.ERROR, "MainActivity", "jmuhia, fetchGreeting, e: ${e}, e.status: ${e.status}")
                Toast.makeText(
                    this@MainActivity,
                    "RUH-ROH! ${e.status.code}",
                    Toast.LENGTH_LONG
                ).show()
            }
        }
    }
}

@Composable
fun Greeting(onClick: () -> Unit) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(vertical = Dp(24f)),
        verticalArrangement = Arrangement.Top,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "Welcome",
        )
        Button(
            onClick = {
                onClick.invoke()
            }
        ) {
            Text(
                text = "Click me"
            )
        }
    }
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    HelloWorldTheme {
        Greeting {}
    }
}
