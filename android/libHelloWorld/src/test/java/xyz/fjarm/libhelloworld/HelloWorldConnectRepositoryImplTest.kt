package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldResponse
import com.connectrpc.BidirectionalStreamInterface
import com.connectrpc.ClientOnlyStreamInterface
import com.connectrpc.Headers
import com.connectrpc.MethodSpec
import com.connectrpc.ProtocolClientInterface
import com.connectrpc.ResponseMessage
import com.connectrpc.ServerOnlyStreamInterface
import com.connectrpc.UnaryBlockingCall
import com.connectrpc.http.Cancelable
import org.junit.Assert.assertThrows
import org.junit.Before
import org.junit.Test

class HelloWorldConnectRepositoryImplTest {

    private lateinit var repository: HelloWorldConnectRepositoryImpl

    @Before
    fun setUp() {
        val fakeClient = object : ProtocolClientInterface {
            override fun <T : Any, U : Any> unary(
                request: T,
                headers: Headers,
                methodSpec: MethodSpec<T, U>,
                onResult: (ResponseMessage<U>) -> Unit
            ): Cancelable = throw UnsupportedOperationException()

            override fun <Input : Any, Output : Any> unaryBlocking(
                request: Input,
                headers: Headers,
                methodSpec: MethodSpec<Input, Output>
            ): UnaryBlockingCall<Output> = throw UnsupportedOperationException()

            override suspend fun <Input : Any, Output : Any> clientStream(
                headers: Headers,
                methodSpec: MethodSpec<Input, Output>
            ): ClientOnlyStreamInterface<Input, Output> = throw UnsupportedOperationException()

            override suspend fun <Input : Any, Output : Any> serverStream(
                headers: Headers,
                methodSpec: MethodSpec<Input, Output>
            ): ServerOnlyStreamInterface<Input, Output> = throw UnsupportedOperationException()

            override suspend fun <Input : Any, Output : Any> stream(
                headers: Headers,
                methodSpec: MethodSpec<Input, Output>
            ): BidirectionalStreamInterface<Input, Output> = throw UnsupportedOperationException()

            override suspend fun <T : Any, U : Any> unary(
                request: T,
                headers: Headers,
                methodSpec: MethodSpec<T, U>
            ): ResponseMessage<U> = throw UnsupportedOperationException()

        }
        repository = HelloWorldConnectRepositoryImpl(fakeClient)
    }

    @Test
    fun getHelloWorld_onResultCallback_throwsNotImplementedError() {
        // Given a valid request
        val request = GetHelloWorldRequest.getDefaultInstance()
        val headers = emptyMap<String, List<String>>()
        val onResult: (ResponseMessage<GetHelloWorldResponse>) -> Unit = {}

        // When getHelloWorld is called with a callback
        // Then throw NotImplementedError
        assertThrows(NotImplementedError::class.java) {
            repository.getHelloWorld(request, headers, onResult)
        }
    }
}
