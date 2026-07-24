package xyz.fjarm.coroutines

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.StandardTestDispatcher
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.setMain
import org.junit.Assert.assertEquals
import org.junit.Test

@OptIn(ExperimentalCoroutinesApi::class)
class DispatchersModuleTest {

    @Test
    fun provideIoDispatcher_returnsIoDispatcher() {
        // When the IO dispatcher is provided
        val result = DispatchersModule.provideIoDispatcher()

        // Then it is Dispatchers.IO
        assertEquals(Dispatchers.IO, result)
    }

    @Test
    fun provideMainDispatcher_returnsMainDispatcher() {
        // Given a test dispatcher is set as Main
        val testDispatcher = StandardTestDispatcher()
        Dispatchers.setMain(testDispatcher)

        try {
            // When the Main dispatcher is provided
            val result = DispatchersModule.provideMainDispatcher()

            // Then it is Dispatchers.Main
            assertEquals(Dispatchers.Main, result)
        } finally {
            Dispatchers.resetMain()
        }
    }
}
