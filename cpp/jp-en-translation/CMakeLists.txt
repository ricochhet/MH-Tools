cmake_minimum_required(VERSION 3.16)
project(termjiten)

set(CMAKE_CXX_STANDARD 23)


set(PUGIXML_SOURCES
    "${CMAKE_SOURCE_DIR}/dependencies/pugixml/src/pugixml.cpp"
)

set(PUGIXML_HEADERS
    "${CMAKE_SOURCE_DIR}/dependencies/pugixml/src/pugixml.hpp"
    "${CMAKE_SOURCE_DIR}/dependencies/pugixml/src/pugiconfig.hpp"
)

include_directories(
    ${CMAKE_SOURCE_DIR}/src
    ${CMAKE_SOURCE_DIR}/dependencies/pugixml/src
)

file(GLOB MAIN_SOURCES "src/*.cpp" "src/*.h")
file(GLOB LIB_SOURCES "src/include/*.cpp" "src/include/*.h")

add_executable(termjiten
    ${MAIN_SOURCES}
    ${LIB_SOURCES}
    ${PUGIXML_SOURCES}
    ${PUGIXML_HEADERS}
)