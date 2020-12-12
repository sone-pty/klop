#include <thread>
#include <mutex>
#include <memory>
#include <chrono>
#include <iostream>
#include <map>
#include <stack>
#include <vector>
#include <string>
#include <algorithm>

#include "SoneUtil/FileUtil/file_util.h"
#include "SoneUtil/StringUtil/string_util.h"

using namespace std;

volatile int counter = 0;
mutex mtx;

void task_lock()
{
    for (int i = 0; i < 1000; ++i)
    {
        lock_guard<mutex> lm(mtx);
        ++counter;
    }
}

void task_trylock()
{
    int n = 0;
    while(n != 1000)
    {
        if(mtx.try_lock())
        {
            ++counter;
            mtx.unlock();
            ++n;
        }
    }
}

int main()
{
    return 0;
}