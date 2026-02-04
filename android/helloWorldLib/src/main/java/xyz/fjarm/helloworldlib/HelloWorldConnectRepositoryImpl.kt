package xyz.fjarm.helloworldlib

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
import java.util.UUID
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class HelloWorldConnectRepositoryImpl @Inject constructor(
    private val client: ProtocolClientInterface,
) : HelloWorldRepository,
    HelloWorldServiceClientInterface
{

    companion object {
        // No need to inject this path. Service and schema upgrades can be done with rolling deploys
        // on backend and force upgrade prompts on the app.
        private const val GET_HELLO_WORLD_V1_PATH =
            "fjarm.helloworld.v1.HelloWorldService/GetHelloWorld"

        private const val CALLBACK_NOT_SUPPORTED_MESSAGE =
            "Connect client does not support callbacks"

        // In a standard repository, the request-id header would be injected using a request header
        // provider.
        private const val REQUEST_ID_HEADER = "request-id"
    }

    override suspend fun getHelloWorld(
        request: GetHelloWorldRequest
    ): GetHelloWorldResponse {
        val headers = mapOf<String, List<String>>(
            REQUEST_ID_HEADER to listOf(UUID.randomUUID().toString())
        )
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
                GET_HELLO_WORLD_V1_PATH,
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
