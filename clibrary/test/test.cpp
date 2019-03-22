#include "librai_tracer.h" // for high_resolution_clock
#include <chrono>          // for high_resolution_clock

static int N = 100;

int main(int argc, char **argv) {
  auto start = std::chrono::high_resolution_clock::now();
  for (int ii = 0; ii < N; ii++) {
    extern GoUintptr SpanStart(GoInt32 p0, GoString p1);
    extern void SpanFinish(GoUintptr p0);
  }
  auto finish = std::chrono::high_resolution_clock::now();
  printf("avg_time = %f\n", (start - finish) / N);
  return 0;
}
