package xyz.fjarm.fjarm

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavKey
import androidx.navigation3.runtime.entryProvider
import androidx.navigation3.ui.NavDisplay
import dagger.hilt.android.AndroidEntryPoint
import xyz.fjarm.fjarmthemelib.FjarmTheme
import xyz.fjarm.navigationlib.NavigatorSideEffect
import xyz.fjarm.navigationlib.Navigator
import javax.inject.Inject

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    @Inject
    lateinit var navigator: Navigator

    @Inject
    lateinit var entryBuilders: Set<@JvmSuppressWildcards EntryProviderScope<NavKey>.() -> Unit>


    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge()
        super.onCreate(savedInstanceState)
        setContent {
            FjarmTheme {
                Scaffold { innerPadding ->
                    Surface(
                        modifier = Modifier
                            .fillMaxSize()
                            .padding(innerPadding),
                    ) {
                        NavDisplay(
                            backStack = navigator.getBackStack(),
                            onBack = {
                                navigator.processSideEffect(
                                    NavigatorSideEffect.NavigateBack,
                                )
                            },
                            entryProvider = entryProvider {
                                entryBuilders.forEach { it() }
                            }
                        )
                    }
                }
            }
        }
    }
}
