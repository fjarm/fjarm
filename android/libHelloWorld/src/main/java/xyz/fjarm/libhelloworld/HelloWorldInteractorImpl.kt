package xyz.fjarm.libhelloworld

import build.buf.gen.fjarm.helloworld.v1.GetHelloWorldRequest
import build.buf.gen.fjarm.helloworld.v1.HelloWorldOutput
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext

class HelloWorldInteractorImpl constructor(
    private val repository: HelloWorldRepository,
) : HelloWorldInteractor {

    companion object {
        private const val ERROR_GREETING = "RUH-ROH!"
    }

    override suspend fun getHelloWorld(): HelloWorldOutput {
        val output = withContext(Dispatchers.IO) {
            try {
                val request = GetHelloWorldRequest.newBuilder().build()
                val response = repository.getHelloWorld(request)
                response.output
            } catch (e: Exception) {
                HelloWorldOutput
                    .newBuilder()
                    .setOutput("$ERROR_GREETING ${e.message}")
                    .build()
            }
        }
        return output
    }
}
