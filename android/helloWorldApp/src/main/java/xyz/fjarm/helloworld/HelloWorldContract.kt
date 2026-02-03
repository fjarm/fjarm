package xyz.fjarm.helloworld

import androidx.annotation.StringRes

data class HelloWorldState(
    @StringRes val promptText: Int,
    @StringRes val buttonText: Int,
)

sealed class HelloWorldEvent {

    data object ButtonClicked: HelloWorldEvent()
}

sealed class HelloWorldSideEffect {

    data class ShowToast(
        val message: String,
    ): HelloWorldSideEffect()
}
