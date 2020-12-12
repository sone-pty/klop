#ifndef SONE_STRINGUTIL_H
#define SONE_STRINGUTIL_H

#include <string>

class StringUtil
{
public:
    static std::wstring StringToWString(const std::string& s);
    static std::string WStringToString(const std::wstring& ws);
};

#endif