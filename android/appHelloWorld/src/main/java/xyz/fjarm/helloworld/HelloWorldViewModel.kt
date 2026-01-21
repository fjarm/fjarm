package xyz.fjarm.helloworld

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableSharedFlow
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asSharedFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import xyz.fjarm.libhelloworld.GetHelloWorldUseCase

class HelloWorldViewModel(
    private val getHelloWorldUseCase: GetHelloWorldUseCase,
): ViewModel() {

    private val _state = MutableStateFlow<HelloWorldState>(HelloWorldState(
        promptText = R.string.prompt_text,
        buttonText = R.string.button_text
    ))
    val state = _state.asStateFlow()

    private val _sideEffect = MutableSharedFlow<HelloWorldSideEffect>()
    val sideEffect = _sideEffect.asSharedFlow()

    fun processEvent(event: HelloWorldEvent) {
        when (event) {
            is HelloWorldEvent.ButtonClicked -> {
                viewModelScope.launch {
                    val helloWorldOutput = getHelloWorldUseCase()
                    _sideEffect.emit(
                        HelloWorldSideEffect.ShowToast(helloWorldOutput.output)
                    )
                }
            }
        }
    }
}
