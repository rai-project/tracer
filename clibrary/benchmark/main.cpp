#include "librai_tracer.h" // for high_resolution_clock
#include <chrono>          // for high_resolution_clock

#include <benchmark/benchmark.h>
#include <string>
#include <thread>

static const GoInt32 NO_TRACE = 0;
static const GoInt32 APPLICATION_TRACE = 1;
static const GoInt32 MODEL_TRACE = 2;
static const GoInt32 FRAMEWORK_TRACE = 3;
static const GoInt32 LIBRARY_TRACE = 4;
static const GoInt32 HARDWARE_TRACE = 5;
static const GoInt32 FULL_TRACE = 6;

namespace detail {

GoString to_go_string(const char *e) {
  GoString res;
  const std::string str = std::string(e);
  res.p = strdup(str.c_str());
  res.n = str.length();
  return res;
}
} // namespace detail

static void BM_Tracer(benchmark::State &state) {
  for (auto _ : state) {
    auto __span = SpanStart(LIBRARY_TRACE, detail::to_go_string("test_run"));
    // std::this_thread::sleep_for(std::chrono::seconds(1));
    benchmark::DoNotOptimize(__span);
    SpanFinish(__span);
  }
}
// Register the function as a benchmark
BENCHMARK(BM_Tracer);

int main(int argc, char **argv) {
  TracerInit();
  benchmark::Initialize(&argc, argv);
  benchmark::RunSpecifiedBenchmarks();
}
