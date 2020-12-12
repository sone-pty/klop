#include "file_util.h"
#include "../StringUtil/string_util.h"

#include <stdio.h>
#include <io.h>
#include <stdlib.h>

FileUtilPtr FileUtil::s_instance = nullptr;
std::mutex FileUtil::s_mutex;

FileUtilPtr FileUtil::GetInstance()
{
    if (s_instance == nullptr)
    {
        std::lock_guard<std::mutex> l(s_mutex);
        if (s_instance == nullptr)
        {
            s_instance = new FileUtil;
        }
    }
    return s_instance;
}

FileUtil::FileUtil()
    :_cache(nullptr)
{
    
}

bool FileUtil::CreateFile(const std::wstring& path)
{
    if(exist(path)) return true;
}

bool FileUtil::exist(const std::wstring& path)
{
    return _access(StringUtil::WStringToString(path).c_str(), 0) == 0;
}