#include "librai_tracer.h" // for high_resolution_clock
#include <chrono>          // for high_resolution_clock

#include <benchmark/benchmark.h>
#include <iostream>
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

static void BM_CTracer(benchmark::State &state) {
  for (auto _ : state) {
    auto iter_span =
        SpanStart(LIBRARY_TRACE, detail::to_go_string("iteration"));
    benchmark::DoNotOptimize(iter_span);
    SpanFinish(iter_span);
  }
}

BENCHMARK(BM_CTracer);

static void BM_CTracerWithContext(benchmark::State &state) {
  SpanStartFromContext_return spanctx =
      SpanStartFromContext(ContextNewBackground(), LIBRARY_TRACE,
                           detail::to_go_string("CTracerWithContext"));
  auto span = spanctx.r0;
  auto ctx = spanctx.r1;
  for (auto _ : state) {
    SpanStartFromContext_return iter_spanctx = SpanStartFromContext(
        ctx, LIBRARY_TRACE, detail::to_go_string("iteration_ctx"));
    auto iter_span = iter_spanctx.r0;
    auto iter_ctx = iter_spanctx.r1;
    std::cout << "ctx = " << iter_ctx << "\n";
    // std::this_thread::sleep_for(std::chrono::seconds(1));
    benchmark::DoNotOptimize(iter_span);
    // benchmark::DoNotOptimize(iter_ctx);
    SpanFinish(iter_span);
    // ContextDelete(iter_ctx);
  }
  SpanFinish(span);
  ContextDelete(ctx);
}

// BENCHMARK(BM_CTracerWithContext);

int main(int argc, char **argv) {
  TracerInit();
  benchmark::Initialize(&argc, argv);
  benchmark::RunSpecifiedBenchmarks();
}