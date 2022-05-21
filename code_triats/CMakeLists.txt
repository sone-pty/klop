# cmake version
cmake_minimum_required(VERSION 3.14)

# project name
project(sone)

# src
aux_source_directory(. CURRENT_SRC_DIR)
aux_source_directory(./extern/http_parser EXTERN_HTTPPARSER_SRC_DIR)
aux_source_directory(./src/ds SRC_DS_DIR)

include_directories(./extern)
include_directories(./src)

# extern lib
set(EXTERNLIB pthread)

# private lib

# generate
add_executable(${CMAKE_PROJECT_NAME} ${CURRENT_SRC_DIR} ${EXTERN_HTTPPARSER_SRC_DIR} ${SRC_DS_DIR})

# link
target_link_libraries(${CMAKE_PROJECT_NAME} ${EXTERNLIB})