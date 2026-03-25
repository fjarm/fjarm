package xyz.fjarm.fjarm

import android.annotation.SuppressLint
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.animation.slideInHorizontally
import androidx.compose.animation.slideOutHorizontally
import androidx.compose.animation.togetherWith
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Surface
import androidx.compose.ui.Modifier
import androidx.navigation3.runtime.EntryProviderScope
import androidx.navigation3.runtime.NavEntry
import androidx.navigation3.runtime.NavKey
import androidx.navigation3.runtime.entryProvider
import androidx.navigation3.ui.NavDisplay
import dagger.hilt.android.AndroidEntryPoint
import xyz.fjarm.fjarmtheme.FjarmTheme
import xyz.fjarm.navigation.Navigation
import xyz.fjarm.navigation.NavigationSideEffect
import xyz.fjarm.underconstructionimpl.UnderConstructionScreen
import javax.inject.Inject

@AndroidEntryPoint
class FjarmActivity : ComponentActivity() {

    @Inject
    lateinit var navigator: Navigation

    @Inject
    lateinit var entryBuilders: Set<@JvmSuppressWildcards EntryProviderScope<NavKey>.() -> Unit>

    @SuppressLint("UnusedMaterial3ScaffoldPaddingParameter")
    override fun onCreate(savedInstanceState: Bundle?) {
        enableEdgeToEdge()
        super.onCreate(savedInstanceState)
        setContent {
            FjarmTheme {
                Scaffold(
                    modifier = Modifier.fillMaxSize(),
                ) { innerPadding ->
                    Surface(
                        modifier = Modifier
                            .fillMaxSize(),
                    ) {
                        NavDisplay(
                            backStack = navigator.getBackStack(),
                            onBack = {
                                navigator.processSideEffect(
                                    NavigationSideEffect.NavigateBack,
                                )
                            },
                            entryProvider = entryProvider(fallback = {
                                NavEntry(it) {
                                    UnderConstructionScreen()
                                }
                            }) {
                                entryBuilders.forEach { it() }
                            },
                            transitionSpec = {
                                // Use togetherWith to create a ContentTransform
                                (slideInHorizontally { fullWidth -> fullWidth } + fadeIn()) togetherWith
                                        fadeOut()
                            },
                            popTransitionSpec = {
                                // For popping, we slide out the exiting screen
                                fadeIn() togetherWith
                                        (slideOutHorizontally { fullWidth -> fullWidth } + fadeOut())
                            },
                            predictivePopTransitionSpec = {
                                // For popping, we slide out the exiting screen
                                fadeIn() togetherWith
                                        (slideOutHorizontally { fullWidth -> fullWidth } + fadeOut())
                            },
                        )
                    }
                }
            }
        }
    }
}
