package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldResponse
import build.buf.gen.fjarm.helloworld.v1.HelloWorldServiceClientInterface
import com.connectrpc.Headers
import com.connectrpc.MethodSpec
import com.connectrpc.ProtocolClientInterface
import com.connectrpc.ResponseMessage
import com.connectrpc.StreamType
import com.connectrpc.getOrThrow
import com.connectrpc.http.Cancelable

class HelloWorldConnectRepositoryImpl(
    private val client: ProtocolClientInterface,
) : HelloWorldRepository,
    HelloWorldServiceClientInterface
{

    companion object {
        // No need to inject this path. Service and schema upgrades can be done with rolling deploys
        // on backend and force upgrade prompts on the app.
        private const val GET_HELLO_WORLD_PATH =
            "fjarm.helloworld.v1.HelloWorldService/GetHelloWorld"

        private const val CALLBACK_NOT_SUPPORTED_MESSAGE =
            "Connect client does not support callbacks"
    }

    override suspend fun getHelloWorld(
        request: GetHelloWorldRequest
    ): GetHelloWorldResponse {
        val headers = emptyMap<String, List<String>>()
        val response = getHelloWorld(request, headers)
        // The fold method can be used for more granular response/error handling.
        return response.getOrThrow()
    }

    override suspend fun getHelloWorld(
        request: GetHelloWorldRequest,
        headers: Headers
    ): ResponseMessage<GetHelloWorldResponse> {
        return client.unary(
            request,
            headers,
            MethodSpec(
                GET_HELLO_WORLD_PATH,
                GetHelloWorldRequest::class,
                GetHelloWorldResponse::class,
                StreamType.UNARY,
            )
        )
    }

    override fun getHelloWorld(
        request: GetHelloWorldRequest,
        headers: Headers,
        onResult: (ResponseMessage<GetHelloWorldResponse>) -> Unit
    ): Cancelable {
        throw NotImplementedError(CALLBACK_NOT_SUPPORTED_MESSAGE)
    }
}
