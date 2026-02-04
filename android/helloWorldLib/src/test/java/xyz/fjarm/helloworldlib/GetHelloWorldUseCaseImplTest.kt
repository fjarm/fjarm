package xyz.fjarm.helloworldlib

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldResponse
import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput
import com.connectrpc.Code
import com.connectrpc.ConnectException
import kotlinx.coroutines.runBlocking
import org.junit.Assert.assertTrue
import org.junit.Test

class GetHelloWorldUseCaseImplTest {

    @Test
    fun invoke_withSuccessAndPresentOutputResponse_returnsHelloWorldOutput() {
        // Given a repository that returns a healthy, fully populated response
        val fakeRepository = object : HelloWorldRepository {
            override suspend fun getHelloWorld(request: GetHelloWorldRequest): GetHelloWorldResponse {
                return GetHelloWorldResponse
                    .newBuilder()
                    .setOutput(
                        HelloWorldOutput
                            .newBuilder()
                            .setOutput("Blah"))
                    .build()
            }
        }
        val getHelloWorldUseCase = GetHelloWorldUseCaseImpl(fakeRepository)

        // When invoking the use case
        val output = runBlocking { getHelloWorldUseCase() }

        // Then no exception is thrown and the output should contain a non-empty string
        assertTrue(output.output.isNotEmpty())
    }

    @Test
    fun invoke_withSuccessAndEmptyOutputResponse_returnsHelloWorldOutput() {
        // Given a repository that returns a healthy, but empty response
        val fakeRepository = object : HelloWorldRepository {
            override suspend fun getHelloWorld(request: GetHelloWorldRequest): GetHelloWorldResponse {
                return GetHelloWorldResponse
                    .newBuilder()
                    .clearOutput()
                    .build()
            }
        }
        val getHelloWorldUseCase = GetHelloWorldUseCaseImpl(fakeRepository)

        // When invoking the use case
        val output = runBlocking { getHelloWorldUseCase() }

        // Then the exception handling in the use case ensures a valid HelloWorldOutput is still
        // returned
        assertTrue(output.output.isNotEmpty())
    }

    @Test
    fun invoke_withErrorResponse_returnsHelloWorldOutput() {
        // Given a repository that returns a error response
        val fakeRepository = object : HelloWorldRepository {
            override suspend fun getHelloWorld(request: GetHelloWorldRequest): GetHelloWorldResponse {
                throw ConnectException(Code.UNKNOWN)
            }
        }
        val getHelloWorldUseCase = GetHelloWorldUseCaseImpl(fakeRepository)

        // When invoking the use case
        val output = runBlocking { getHelloWorldUseCase() }

        // Then the exception handling in the use case ensures a valid HelloWorldOutput is still
        // returned
        assertTrue(output.output.isNotEmpty())
    }
}
