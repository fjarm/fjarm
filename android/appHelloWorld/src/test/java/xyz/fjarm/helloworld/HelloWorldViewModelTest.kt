package xyz.fjarm.helloworld

import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.launch
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.Assert.assertEquals
import org.junit.Before
import org.junit.Test
import xyz.fjarm.libhelloworld.GetHelloWorldUseCase

@OptIn(ExperimentalCoroutinesApi::class)
class HelloWorldViewModelTest {

    private val testDispatcher = StandardTestDispatcher()

    @Before
    fun setUp() {
        Dispatchers.setMain(testDispatcher)
    }

    @Test
    fun processEvent_ButtonClicked_withSuccessUseCaseInvoke_emitsShowToastSideEffect() = runTest {
        // Given a use case that returns valid, non empty output
        val expectedOutput = "blah blah"
        val fakeUseCase = object : GetHelloWorldUseCase {
            override suspend fun invoke(): HelloWorldOutput {
                return HelloWorldOutput.newBuilder()
                    .setOutput(expectedOutput)
                    .build()
            }
        }
        val viewModel = HelloWorldViewModel(fakeUseCase)
        
        val collectedSideEffects = mutableListOf<HelloWorldSideEffect>()
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.sideEffect.collect { collectedSideEffects.add(it) }
        }

        // When a ButtonClicked event is processed by the ViewModel
        viewModel.processEvent(HelloWorldEvent.ButtonClicked)
        advanceUntilIdle()

        // Then a ShowToast side effect is emitted with the expected output
        val sideEffect = collectedSideEffects.firstOrNull()
        assert(collectedSideEffects.size == 1)
        assert(sideEffect is HelloWorldSideEffect.ShowToast)
        assertEquals(
            expectedOutput,
            (sideEffect as HelloWorldSideEffect.ShowToast).message,
        )
    }

    @Test
    fun processEvent_ButtonClicked_withSuccessUseCaseInvoke_stateRemainsConstant() = runTest {
        // Given a use case that returns valid, non empty output
        val expectedOutput = "blah blah"
        val fakeUseCase = object : GetHelloWorldUseCase {
            override suspend fun invoke(): HelloWorldOutput {
                return HelloWorldOutput.newBuilder()
                    .setOutput(expectedOutput)
                    .build()
            }
        }
        val viewModel = HelloWorldViewModel(fakeUseCase)

        val collectedStates = mutableListOf<HelloWorldState>()
        backgroundScope.launch(UnconfinedTestDispatcher(testScheduler)) {
            viewModel.state.collect { collectedStates.add(it) }
        }

        // When a ButtonClicked event is processed by the ViewModel
        viewModel.processEvent(HelloWorldEvent.ButtonClicked)
        advanceUntilIdle()

        // Then the state remains constant and does not emit more than once
        val state = collectedStates.firstOrNull()
        assert(collectedStates.size == 1)
        assert(state is HelloWorldState)
        assertEquals(R.string.prompt_text, (state as HelloWorldState).promptText)
        assertEquals(R.string.button_text, state.buttonText)
    }
}
