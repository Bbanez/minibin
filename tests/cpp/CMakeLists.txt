cmake_minimum_required(VERSION 4.1)
project(minibin_test VERSION 0.1.0)
add_executable(minibin_test main.cpp minibin.cpp)

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_CXX_STANDARD 20)
