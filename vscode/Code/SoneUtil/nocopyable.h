#ifndef SONE_NOCOPYABLE_H
#define SONE_NOCOPYABLE_H

class NoCopyable
{
public:
    NoCopyable(const NoCopyable&) = delete;
    void operator=(const NoCopyable&) = delete;

protected:
    NoCopyable() = default;
    ~NoCopyable() = default;
};

#endif