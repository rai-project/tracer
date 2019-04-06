cmake_minimum_required(VERSION 3.1)


include("cmake/HunterGate.cmake")
HunterGate(
    URL "https://github.com/ruslo/hunter/archive/v0.23.142.tar.gz"
    SHA1 "1d841ca74150c92054896a7f6d88a123dd9e615d"
)


project(tracer-benchmark LANGUAGES CXX)

hunter_add_package(benchmark)

find_package(benchmark CONFIG REQUIRED)

if (CMAKE_SYSTEM_NAME STREQUAL "Linux")
	set(SYSTEMID "Linux-x86-64")
elseif (CMAKE_SYSTEM_NAME STREQUAL "Darwin")
	set(SYSTEMID "MacOSX-x86-64")
elseif (CMAKE_SYSTEM_NAME STREQUAL "Windows")
	set(SYSTEMID "Windows-x86-64")
endif()

set(CLibrary_dist_dir ${PROJECT_SOURCE_DIR}/../dist/${SYSTEMID})

SET(CMAKE_FIND_LIBRARY_SUFFIXES .a .so .dylib)

add_executable(bench main.cpp)
set_property(TARGET bench PROPERTY CXX_STANDARD 14)
target_include_directories(bench PUBLIC ${CLibrary_dist_dir})
find_library(TRACER_LIBRARY rai_tracer HINTS ${CLibrary_dist_dir})
message(STATUS "lib = ${TRACER_LIBRARY}")
target_link_libraries(bench benchmark::benchmark ${TRACER_LIBRARY} pthread)
if(SYSTEMID STREQUAL "Linux-x86-64")
	find_library(NVML_LIBRARY nvidia-ml HINTS /usr/lib/x86_64-linux-gnu)
	target_link_libraries(bench benchmark::benchmark ${NVML_LIBRARY})
endif()
if(APPLE)
  set(CMAKE_FIND_FRAMEWORK ONLY)
  find_library(CoreFoundation_LIBRARY CoreFoundation)
  find_library(Security_LIBRARY Security)
  target_link_libraries(bench ${CoreFoundation_LIBRARY} ${Security_LIBRARY})
endif(APPLE)
