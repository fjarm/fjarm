package xyz.fjarm.helloworld

import android.os.Bundle
import android.widget.Toast
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.Dp
import androidx.hilt.lifecycle.viewmodel.compose.hiltViewModel
import androidx.lifecycle.Lifecycle
import androidx.lifecycle.compose.LocalLifecycleOwner
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.lifecycle.repeatOnLifecycle
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.StateFlow
import xyz.fjarm.helloworld.ui.theme.HelloWorldTheme

@AndroidEntryPoint
class MainActivity: ComponentActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge()
        super.onCreate(savedInstanceState)
        setContent {
            HelloWorldTheme {
                // Scaffold is a layout structure for the entire screen.
                Scaffold { innerPadding ->
                    // A surface container using the 'background' color from the theme
                    Surface(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(innerPadding),
                        color = MaterialTheme.colorScheme.background,
                    ) {
                        Greeting()
                    }
                }
            }
        }
    }
}

@Composable
fun Greeting(
    modifier: Modifier = Modifier,
    helloWorldViewModel: HelloWorldViewModel = hiltViewModel<HelloWorldViewModel>(),
) {
    Greeting(
        modifier = modifier,
        state = helloWorldViewModel.state,
        sideEffect = helloWorldViewModel.sideEffect,
        processEvent = helloWorldViewModel::processEvent,
    )
}

@Composable
fun Greeting(
    modifier: Modifier = Modifier,
    state: StateFlow<HelloWorldState>,
    sideEffect: SharedFlow<HelloWorldSideEffect>,
    processEvent: (HelloWorldEvent) -> Unit,
) {
    val context = LocalContext.current
    val lifecycleOwner = LocalLifecycleOwner.current

    // If the SharedFlow changes, cancel the old collection and start a new one
    LaunchedEffect(sideEffect, lifecycleOwner) {
        // Ensure that side effect collection only happens when the UI is in the STARTED state.
        lifecycleOwner.repeatOnLifecycle(Lifecycle.State.STARTED) {
            sideEffect.collect { sideEffect ->
                when (sideEffect) {
                    is HelloWorldSideEffect.ShowToast -> {
                        Toast.makeText(
                            context,
                            sideEffect.message,
                            Toast.LENGTH_SHORT,
                        ).show()
                    }
                }
            }
        }
    }

    val state = state.collectAsStateWithLifecycle()

    Column(
        modifier = Modifier
            .padding(Dp(24f)),
        verticalArrangement = Arrangement.Top,
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(text = stringResource(id = state.value.promptText))
        Button(
            onClick = { processEvent(HelloWorldEvent.ButtonClicked) }
        ) {
            Text(text = stringResource(id = state.value.buttonText))
        }
    }
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    HelloWorldTheme {
        Greeting(
            modifier = Modifier,
            state = MutableStateFlow(
                HelloWorldState(
                    promptText = R.string.prompt_text,
                    buttonText = R.string.button_text,
                )
            ),
            sideEffect = MutableSharedFlow(),
            processEvent = { },
        )
    }
}
