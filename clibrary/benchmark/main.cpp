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

static void BM_CTracer(benchmark::State &state) {
  for (auto _ : state) {
    auto iter_span = SpanStart(APPLICATION_TRACE, (char *)"iteration");
    benchmark::DoNotOptimize(iter_span);
    SpanFinish(iter_span);
  }
}

// BENCHMARK(BM_CTracer);

static void BM_CTracerWithContext(benchmark::State &state) {
  SpanStartFromContext_return spanctx = SpanStartFromContext(
      ContextNewBackground(), APPLICATION_TRACE, (char *)"CTracerWithContext");
  auto span = spanctx.r0;
  auto ctx = spanctx.r1;
  int ii = 0;
  for (auto _ : state) {
    if (ii++ > 300) {
      state.SkipWithError("limit");
      break;
    }
    SpanStartFromContext_return iter_spanctx =
        SpanStartFromContext(ctx, APPLICATION_TRACE, (char *)"iteration_ctx");
    auto iter_span = iter_spanctx.r0;
    auto iter_ctx = iter_spanctx.r1;
    benchmark::DoNotOptimize(iter_span);
    benchmark::DoNotOptimize(iter_ctx);
    SpanFinish(iter_span);
    ContextDelete(iter_ctx);
  }
  SpanFinish(span);
  ContextDelete(ctx);
}

BENCHMARK(BM_CTracerWithContext);

void test() {

  SpanStartFromContext_return spanctx1 = SpanStartFromContext(
      ContextNewBackground(), APPLICATION_TRACE, (char *)"CTracerWithContext");
  auto span1 = spanctx1.r0;
  auto ctx1 = spanctx1.r1;

  SpanStartFromContext_return spanctx2 =
      SpanStartFromContext(ctx1, APPLICATION_TRACE, (char *)"iteration");
  auto span2 = spanctx2.r0;
  auto ctx2 = spanctx2.r1;

  std::cout << "Ctx1 = " << std::hex << ctx1 << "\n";
  std::cout << "Ctx2 = " << std::hex << ctx2 << "\n";

  // ContextDelete(ctx2);
  SpanFinish(span2);

  // ContextDelete(ctx1);
  SpanFinish(span1);
}

int main(int argc, char **argv) {
  TracerInit();
  TracerSetLevel(FULL_TRACE);
  benchmark::Initialize(&argc, argv);
  benchmark::RunSpecifiedBenchmarks();
  // test();

  // std::this_thread::sleep_for(std::chrono::milliseconds(1000));
  TracerClose();
}
