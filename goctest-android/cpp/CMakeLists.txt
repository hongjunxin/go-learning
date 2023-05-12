cmake_minimum_required(VERSION 3.0.0)
project(hello VERSION 0.0.1)

add_executable(${PROJECT_NAME} hello.cpp)

set_target_properties(${PROJECT_NAME} PROPERTIES
    CXX_STANDARD 17
    OUTPUT_NAME ${PROJECT_NAME})

target_include_directories(${PROJECT_NAME} PRIVATE ../go)
target_link_directories(${PROJECT_NAME} PRIVATE ../go)
target_link_libraries(${PROJECT_NAME} goc log)
