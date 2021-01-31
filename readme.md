Shows memory usage of a simple Ebiten game with a 4K (3840 x 2160) layout.

Can be viewed here: https://vzx.github.io/ebiten-memory-usage/

The stats displayed are:

| Item | Explanation |
| --- | --- |
| `vp` | The viewpoint `x16` and `y16`. |
| `ticks` | Total number of times the `Update()` function has been called. |
| `Alloc` | The current number of allocated bytes. |
| `Total` | The total number of bytes that have been allocated (and potentially garbage collected) so far. |
| `Sys` | The number of bytes the runtime has obtained from the operating system. |
| `NextGC` | The target heap size of the next GC cycle. |
| `NumGC` | Number of garbage collection cycles run so far. |

See also `runtime.ReadMemStats()`: https://golang.org/pkg/runtime/#ReadMemStats
