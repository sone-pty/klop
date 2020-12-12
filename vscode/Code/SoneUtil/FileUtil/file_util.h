#ifndef SONE_FILEUTIL_H
#define SONE_FILEUTIL_H

#include "../nocopyable.h"

#include <mutex>
#include <string>

class FileUtil;

typedef FileUtil* FileUtilPtr;

class FileUtil : public NoCopyable
{
public:
    static FileUtilPtr GetInstance();

public:
    bool CreateFile(const std::wstring& path);
    bool IsExist(const std::wstring& path) { return exist(path); }

protected:
    FileUtil();
    FileUtil(const FileUtil&) = delete;
    FileUtil operator=(const FileUtil&) = delete;

private:
    bool exist(const std::wstring& path);

private:
    FILE* _cache;   
    static FileUtilPtr s_instance;
    static std::mutex s_mutex;
};

#endif