package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput
import com.connectrpc.Code
import com.connectrpc.ConnectException
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class GetHelloWorldUseCaseImpl @Inject constructor(
    private val repository: HelloWorldRepository,
) : GetHelloWorldUseCase {

    companion object {
        private const val ERROR_GREETING = "RUH-ROH!"
        private const val HIDDEN_ERROR = "Something went wrong. Try again later."
        private const val UNAVAILABLE_MESSAGE = "Something went wrong on our end. Try again later."
    }

    /**
     * This particular use case does not return a [Result] because the same behavior applies in the
     * success and failure case. A toast with a message is shown in both cases.
     */
    override suspend fun invoke(): HelloWorldOutput {
        return withContext(Dispatchers.IO) {
            try {
                val request = GetHelloWorldRequest.newBuilder().build()
                val response = repository.getHelloWorld(request)
                response.output
            } catch (e: ConnectException) {
                // ConnectRPC always throws a ConnectException. Specific error handling has to look
                // at the contained error code.
                val message = when (e.code) {
                    Code.UNAVAILABLE -> UNAVAILABLE_MESSAGE
                    else -> {
                        // Consider logging this unexpected exception
                        HIDDEN_ERROR
                    }
                }
                HelloWorldOutput
                    .newBuilder()
                    .setOutput("$ERROR_GREETING $message}")
                    .build()
            } catch (_: Exception) {
                // Generic error handling is reserved for truly exceptional errors that happen
                // outside of ConnectRPC.
                HelloWorldOutput
                    .newBuilder()
                    .setOutput("$ERROR_GREETING $HIDDEN_ERROR")
                    .build()
            }
        }
    }
}
