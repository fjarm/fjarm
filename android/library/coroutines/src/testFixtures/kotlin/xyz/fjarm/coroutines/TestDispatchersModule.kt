package xyz.fjarm.coroutines

import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import kotlinx.coroutines.CoroutineDispatcher
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.UnconfinedTestDispatcher

/**
 * A Hilt module that provides [UnconfinedTestDispatcher] for both [IoDispatcher] and [MainDispatcher].
 *
 * This can be used in tests to override the production [DispatchersModule].
 *
 * Example usage in an Android test:
 * ```
 * @Module
 * @TestInstallIn(
 *     components = [SingletonComponent::class],
 *     replaces = [DispatchersModule::class]
 * )
 * object TestDispatchersModule {
 *     @Provides @IoDispatcher fun provideIoDispatcher(): CoroutineDispatcher = UnconfinedTestDispatcher()
 *     @Provides @MainDispatcher fun provideMainDispatcher(): CoroutineDispatcher = UnconfinedTestDispatcher()
 * }
 * ```
 *
 * Note: Since `@TestInstallIn` is part of Hilt's Android testing library, it is not used here
 * to keep this module pure Kotlin. Consumers in Android modules can use it to replace
 * [DispatchersModule] with these test dispatchers.
 */
@OptIn(ExperimentalCoroutinesApi::class)
@Module
@InstallIn(SingletonComponent::class)
object TestDispatchersModule {

    @Provides
    @IoDispatcher
    fun provideIoDispatcher(): CoroutineDispatcher {
        return UnconfinedTestDispatcher()
    }

    @Provides
    @MainDispatcher
    fun provideMainDispatcher(): CoroutineDispatcher {
        return UnconfinedTestDispatcher()
    }
}
